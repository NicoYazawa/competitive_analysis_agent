package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/lib/pq"
)

// VectorStore pgvector向量存储
type VectorStore struct {
	db     *PostgresDB
	logger *slog.Logger
}

// NewVectorStore 创建向量存储实例
func NewVectorStore(db *PostgresDB, logger *slog.Logger) *VectorStore {
	return &VectorStore{
		db:     db,
		logger: logger,
	}
}

// VectorEmbedding 向量嵌入
type VectorEmbedding struct {
	ID       string    `json:"id"`
	Entity   string    `json:"entity"`   // competitor, product, trend
	Data     string    `json:"data"`     // JSON数据
	Vector   []float32 `json:"vector"`   // 1536维向量
	Metadata JSONMap   `json:"metadata"` // 额外元数据
	CreatedAt time.Time `json:"created_at"`
}

// JSONMap JSON映射类型
type JSONMap map[string]interface{}

// InitVectorTable 初始化向量表
func (v *VectorStore) InitVectorTable(ctx context.Context) error {
	// 创建向量表
	schema := `
	CREATE TABLE IF NOT EXISTS embeddings (
		id TEXT PRIMARY KEY,
		entity TEXT NOT NULL,
		data JSONB NOT NULL,
		vector vector(1536) NOT NULL,
		metadata JSONB DEFAULT '{}',
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_embeddings_entity ON embeddings(entity);
	CREATE INDEX IF NOT EXISTS idx_embeddings_vector ON embeddings USING ivfflat(vector vector_cosine_ops);
	`

	_, err := v.db.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("create vector table failed: %w", err)
	}

	v.logger.Info("Vector table initialized")
	return nil
}

// Insert 插入向量
func (v *VectorStore) Insert(ctx context.Context, emb *VectorEmbedding) error {
	metadataJSON, err := json.Marshal(emb.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata failed: %w", err)
	}

	query := `
	INSERT INTO embeddings (id, entity, data, vector, metadata, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE SET
		data = EXCLUDED.data,
		vector = EXCLUDED.vector,
		metadata = EXCLUDED.metadata
	`

	_, err = v.db.Exec(ctx, query,
		emb.ID,
		emb.Entity,
		emb.Data,
		pq.Array(emb.Vector),
		metadataJSON,
		emb.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert embedding failed: %w", err)
	}

	return nil
}

// SearchTopK 向量相似度搜索TopK
func (v *VectorStore) SearchTopK(ctx context.Context, entity string, queryVector []float32, k int) ([]*VectorEmbedding, error) {
	query := `
	SELECT id, entity, data, vector, metadata, created_at
	FROM embeddings
	WHERE entity = $1
	ORDER BY vector <=> $2
	LIMIT $3
	`

	rows, err := v.db.Query(ctx, query, entity, pq.Array(queryVector), k)
	if err != nil {
		return nil, fmt.Errorf("search vectors failed: %w", err)
	}
	defer rows.Close()

	var results []*VectorEmbedding
	for rows.Next() {
		var emb VectorEmbedding
		var vectorBytes []byte
		var metadataJSON []byte

		err := rows.Scan(
			&emb.ID,
			&emb.Entity,
			&emb.Data,
			&vectorBytes,
			&metadataJSON,
			&emb.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}

		// 反序列化向量
		if err := json.Unmarshal(vectorBytes, &emb.Vector); err != nil {
			// 尝试直接解析
			emb.Vector = bytesToFloat32(vectorBytes)
		}

		// 反序列化元数据
		if err := json.Unmarshal(metadataJSON, &emb.Metadata); err != nil {
			emb.Metadata = make(JSONMap)
		}

		results = append(results, &emb)
	}

	return results, nil
}

// SearchWithFilter 带过滤条件的向量搜索
func (v *VectorStore) SearchWithFilter(ctx context.Context, entity string, queryVector []float32, k int, filter map[string]interface{}) ([]*VectorEmbedding, error) {
	// 构建过滤条件
	filterConditions := "entity = $1"
	args := []interface{}{entity}
	argIdx := 4

	for key, value := range filter {
		filterConditions += fmt.Sprintf(" AND metadata->>'%s' = $%d", key, argIdx)
		args = append(args, value)
		argIdx++
	}

	query := fmt.Sprintf(`
	SELECT id, entity, data, vector, metadata, created_at
	FROM embeddings
	WHERE %s
	ORDER BY vector <=> $2
	LIMIT $3
	`, filterConditions)

	args = append([]interface{}{}, queryVector, k)
	rows, err := v.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("search vectors with filter failed: %w", err)
	}
	defer rows.Close()

	var results []*VectorEmbedding
	for rows.Next() {
		var emb VectorEmbedding
		var vectorBytes []byte
		var metadataJSON []byte

		err := rows.Scan(
			&emb.ID,
			&emb.Entity,
			&emb.Data,
			&vectorBytes,
			&metadataJSON,
			&emb.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan row failed: %w", err)
		}

		if err := json.Unmarshal(vectorBytes, &emb.Vector); err != nil {
			emb.Vector = bytesToFloat32(vectorBytes)
		}
		if err := json.Unmarshal(metadataJSON, &emb.Metadata); err != nil {
			emb.Metadata = make(JSONMap)
		}

		results = append(results, &emb)
	}

	return results, nil
}

// Delete 删除向量
func (v *VectorStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM embeddings WHERE id = $1`
	_, err := v.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete embedding failed: %w", err)
	}
	return nil
}

// DeleteByEntity 按实体类型删除
func (v *VectorStore) DeleteByEntity(ctx context.Context, entity string) error {
	query := `DELETE FROM embeddings WHERE entity = $1`
	_, err := v.db.Exec(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("delete embeddings by entity failed: %w", err)
	}
	return nil
}

// Count 获取向量数量
func (v *VectorStore) Count(ctx context.Context, entity string) (int, error) {
	query := `SELECT COUNT(*) FROM embeddings WHERE entity = $1`
	var count int
	err := v.db.QueryRow(ctx, query, entity).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count embeddings failed: %w", err)
	}
	return count, nil
}

// CosineSimilarity 计算余弦相似度
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// bytesToFloat32 将字节数组转换为float32数组
func bytesToFloat32(data []byte) []float32 {
	if len(data) == 0 {
		return nil
	}

	// 尝试从JSON解析
	var result []float32
	if err := json.Unmarshal(data, &result); err == nil {
		return result
	}

	return nil
}
