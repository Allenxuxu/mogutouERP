<template>
  <div class="app-container">
    <el-form v-enterToNext="true" ref="postForm" :model="postForm" :rules="rules" status-icon label-width="80px" class="demo-ruleForm">
      <el-form-item label="商品编号" prop="id">
        <el-input v-model="postForm.id"/>
      </el-form-item>
      <el-form-item label="名称" prop="name">
        <el-input v-model="postForm.name"/>
      </el-form-item>
      <el-form-item label="颜色" prop="colour">
        <el-input v-model="postForm.colour"/>
      </el-form-item>
      <el-form-item label="品牌" prop="brand">
        <el-input v-model="postForm.brand"/>
      </el-form-item>
      <el-form-item label="尺寸" prop="size">
        <el-input v-model="postForm.size"/>
      </el-form-item>
      <el-form-item label="售价" prop="price">
        <el-input v-model.number="postForm.price"/>
      </el-form-item>
      <el-form-item label="进价" prop="purchase_price">
        <el-input v-model.number="postForm.purchase_price"/>
      </el-form-item>

      <el-form-item>
        <el-button :loading="loading" type="primary" size="mini" @click="submitForm">提交</el-button>
        <el-button size="mini" @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>

import { createCommodity } from '@/api/commodity'

export default {
  data() {
    return {
      loading: false,
      postForm: {
        id: '',
        name: '',
        colour: '',
        brand: '',
        size: '',
        price: 0,
        purchase_price: 0
      },
      rules: {
        id: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        colour: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        brand: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        size: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        price: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        purchase_price: [
          { required: true, message: '请输入', trigger: 'blur' }
        ]
      }
    }
  },

  methods: {
    submitForm() {
      this.$refs['postForm'].validate(valid => {
        if (valid) {
          this.loading = true
          createCommodity(this.postForm)
            .then(response => {
              console.log(response)
              this.loading = false
              this.resetForm()
              this.$message.success('创建成功')
            })
            .catch(() => {
              this.loading = false
              this.$message.error('创建失败')
            })
        }
      })
    },
    resetForm() {
      this.postForm.id = ''
      this.postForm.name = ''
      this.postForm.colour = ''
      this.postForm.brand = ''
      this.postForm.size = ''
      this.postForm.price = 0
      this.postForm.purchase_price = 0
    }
  }
}
</script>
