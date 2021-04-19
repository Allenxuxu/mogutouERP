<template>
  <div class="app-container">
    <el-table
      v-loading="listLoading"
      :data="data"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row>
      <el-table-column label="姓名" align="center">
        <template slot-scope="scope">
          {{ scope.row.Name }}
        </template>
      </el-table-column>
      <el-table-column label="电话" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.Tel }}</span>
        </template>
      </el-table-column>
      <el-table-column label="职位" align="center">
        <template slot-scope="scope">
          {{ scope.row.Position }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getStaff } from '@/api/company'

export default {
  data() {
    return {
      data: null,
      listLoading: true
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      getStaff().then(response => {
        console.log(response)
        this.data = response
        this.listLoading = false
      })
    }
  }
}
</script>
