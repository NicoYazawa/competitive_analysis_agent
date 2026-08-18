<template>
  <div class="price-table">
    <table>
      <thead>
        <tr>
          <th>Competitor</th>
          <th>Price</th>
          <th>Rating</th>
          <th>Reviews</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in data" :key="item.id">
          <td>{{ item.name }}</td>
          <td>${{ item.price }}</td>
          <td>{{ item.rating }} ★</td>
          <td>{{ item.reviews }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  data: any[]
}>()
</script>

<style scoped>
.price-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}
th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
}
</style>
