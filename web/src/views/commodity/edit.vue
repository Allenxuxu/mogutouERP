<template>
  <div class="app-container">
    <el-input v-model="search" placeholder="输入名称关键字搜索" size="mini" align="right"/>
    <el-table
      v-loading="Loading"
      :data="data.filter(data => !search || data.Name.toLowerCase().includes(search.toLowerCase()))"
      element-loading-text="Loading"
      style="width: 100%"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <el-form label-position="left" inline class="table-expand">
            <el-form-item label="品牌">
              <span>{{ props.row.Brand }}</span>
            </el-form-item>
            <el-form-item label="库存">
              <span>{{ props.row.Number }}</span>
            </el-form-item>
            <el-form-item label="预售量">
              <span>{{ props.row.PresaleNumber }}</span>
            </el-form-item>
            <el-form-item label="销量">
              <span>{{ props.row.SalesVolume }}</span>
            </el-form-item>
            <el-form-item label="售价">
              <span>{{ props.row.Price }}</span>
            </el-form-item>
            <el-form-item label="进价">
              <span>{{ props.row.PurchasePrice }}</span>
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <el-table-column label="ID" prop="ID"/>
      <el-table-column label="名称" prop="Name"/>
      <el-table-column label="颜色" prop="Colour"/>
      <el-table-column label="尺寸" prop="Size"/>

      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button size="mini" type="info" icon="el-icon-edit" circle @click="handleEdit(scope.$index, scope.row)"/>
          <el-button size="mini" type="danger" icon="el-icon-delete" circle @click="handleDelete(scope.$index, scope.row)"/>
        </template>
      </el-table-column>
    </el-table>

    <!-- 跟新表单 -->
    <el-dialog :visible.sync="dialogFormVisible" title="修改商品信息" width="99%">
      <el-form v-enterToNext="true" ref="postForm" :model="postForm" :rules="rules" status-icon label-width="50px" >
        <el-form-item label="ID" prop="ID">
          <span>{{ postForm.ID }}</span>
        </el-form-item>

        <el-form-item label="名称" prop="Name">
          <el-input v-model="postForm.Name"/>
        </el-form-item>
        <el-form-item label="颜色" prop="Colour">
          <el-input v-model="postForm.Colour"/>
        </el-form-item>
        <el-form-item label="品牌" prop="Brand">
          <el-input v-model="postForm.Brand"/>
        </el-form-item>
        <el-form-item label="尺寸" prop="Size">
          <el-input v-model="postForm.Size"/>
        </el-form-item>
        <el-form-item label="售价" prop="Price">
          <el-input v-model.number="postForm.Price"/>
        </el-form-item>
        <el-form-item label="进价" prop="Purchase_Price">
          <el-input v-model.number="postForm.Purchase_Price"/>
        </el-form-item>

      </el-form>
      <div slot="footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="handleEditConfirm">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<style>
.table-expand {
  font-size: 0;
}
.table-expand label {
  width: 90px;
  color: #99a9bf;
}
.table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
</style>

<script>
import { getCommoditiesAsAdmin, deleteCommodity, updateCommodity } from '@/api/commodity'

export default {
  data() {
    return {
      Loading: true,
      data: [],
      search: '',
      dialogFormVisible: false,
      postForm: {
        ID: '',
        Name: '',
        Colour: '',
        Brand: '',
        Size: '',
        Price: 0,
        Purchase_Price: 0
      },
      rules: {
        Name: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Colour: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Brand: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Size: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Price: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Purchase_Price: [
          { required: true, message: '请输入', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.Loading = true
      getCommoditiesAsAdmin().then(response => {
        console.log(response)
        this.data = response
        this.Loading = false
      })
    },
    handleDelete(index, data) {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          this.Loading = true
          console.log(data.ID)
          deleteCommodity(data.ID)
            .then(response => {
              console.log(response)
              this.$message.success('删除成功')
              this.data.splice(index, 1)
              this.Loading = false
            })
            .catch(() => {
              this.Loading = false
              this.$message.error('删除失败')
            })
        })
        .catch(() => {
          this.$message.info('已取消删除')
        })
    },
    handleEdit(index, data) {
      this.dialogFormVisible = true
      this.postForm.ID = data.ID

      this.postForm.Name = data.Name
      this.postForm.Colour = data.Colour
      this.postForm.Brand = data.Brand
      this.postForm.Size = data.Size
      this.postForm.Price = data.Price
      this.postForm.Purchase_Price = data.PurchasePrice
    },
    handleEditConfirm() {
      this.Loading = true
      updateCommodity(this.postForm.ID, this.postForm)
        .then(response => {
          console.log(response)
          this.$message.success('更新成功')
          this.Loading = false
          this.dialogFormVisible = false
          this.fetchData()
        })
        .catch(() => {
          this.Loading = false
          this.$message.error('更新失败')
        })
    }
  }
}
</script>
