<template>
  <div class="app-container">
    <el-form v-enterToNext="true" ref="postForm" :model="postForm" :rules="rules" label-position="left" label-width="110px" size="mini" status-icon >
      <el-form-item label="客户姓名" prop="Name">
        <el-input v-model="postForm.Name" />
      </el-form-item>
      <el-form-item label="客户电话" prop="Tel">
        <el-input v-model="postForm.Tel"/>
      </el-form-item>
      <el-form-item label="预计送货时间" prop="DeliveryTime">
        <el-input v-model="postForm.DeliveryTime"/>
      </el-form-item>
      <el-form-item label="预计送货地址" prop="DeliveryAddress">
        <el-input v-model="postForm.DeliveryAddress"/>
      </el-form-item>

      <el-form-item label="订单总额" prop="Amount">
        <el-input v-model.number="postForm.Amount" type="number"/>
      </el-form-item>
      <el-form-item label="定金" prop="Deposit">
        <el-input v-model.number="postForm.Deposit" type="number"/>
      </el-form-item>
      <el-form-item label="备注" prop="Remarks">
        <el-input v-model="postForm.Remarks"/>
      </el-form-item>

    </el-form>

    <el-select v-model="selectValue" filterable clearable placeholder="请选择" size="mini" @change="selectChange">
      <el-option
        v-for="item in options"
        :key="item.ID"
        :label="item.ID"
        :value="item"
        :disabled="item.Disabled"/>
    </el-select>
    <el-table :data="commodities" style="width: 100%">
      <el-table-column label="ID" prop="ID"/>
      <el-table-column label="名称" prop="Name"/>
      <el-table-column label="颜色" prop="Colour"/>
      <el-table-column label="尺寸" prop="Size"/>
      <el-table-column label="品牌" prop="Brand"/>
      <el-table-column label="数量" width="150">
        <template slot-scope="scope">
          <el-input-number v-model="scope.row.Quantity" :min="0" size="mini"/>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button size="mini" type="danger" icon="el-icon-delete" circle @click="handleDelete(scope.$index,scope.row)"/>
        </template>
      </el-table-column>
    </el-table>
    <br >
    <div>
      <el-button :loading="loading" type="primary" size="mini" @click="submitForm">提交</el-button>
      <el-button size="mini" @click="resetForm">重置</el-button>
    </div>
  </div>
</template>

<script>

import { getCommodities } from '@/api/commodity'
import { createCustormerOrder } from '@/api/order'
import { validateTel } from '@/utils/validate'

export default {
  data() {
    const checkTel = (rule, value, callback) => {
      if (!validateTel(value)) {
        callback(new Error('必须是11位合法手机号'))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      selectValue: '',
      options: [],
      commodities: [],
      postForm: {
        Name: '',
        Tel: '',
        DeliveryAddress: '',
        DeliveryTime: '',
        Amount: null,
        Deposit: null,
        Remarks: '',
        Goods: []
      },
      rules: {
        Name: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Tel: [
          { required: true, validator: checkTel, trigger: 'blur' }
        ],
        DeliveryAddress: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        DeliveryTime: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Deposit: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Amount: [
          { required: true, message: '请输入', trigger: 'blur' }
        ],
        Remarks: [
          { max: 255, message: '最多255  个字符', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getCommodities().then(response => {
        console.log(response)
        this.options = response.map(item => {
          return {
            Quantity: 0,
            ID: item.ID,
            Name: item.Name,
            Colour: item.Colour,
            Size: item.Size,
            Brand: item.Brand
          }
        })
      })
    },
    submitForm() {
      this.$refs['postForm'].validate(valid => {
        if (valid) {
          const len = this.commodities.length
          if (len <= 0) {
            this.$message.error('请添加商品')
            return
          }
          for (let i = 0; i < len; i++) {
            if (this.commodities[i].Quantity <= 0) {
              this.$message.error('商品数量不能为 0')
              this.postForm.Goods = []
              return
            }
            this.postForm.Goods.push({
              ID: this.commodities[i].ID,
              Number: this.commodities[i].Quantity
            })
          }

          console.log(this.postForm)
          this.loading = true
          createCustormerOrder(this.postForm)
            .then(response => {
              console.log(response)
              this.loading = false
              this.resetForm()
              this.$message.success('创建成功')
            })
            .catch(() => {
              this.loading = false
              this.postForm.Goods = []
              this.$message.error('创建失败')
            })
        }
      })
    },
    selectChange() {
      if (this.selectValue !== '') {
        this.selectValue.Disabled = true
        this.commodities.push(this.selectValue)
        this.selectValue = ''
      }
    },
    handleDelete(index, row) {
      this.commodities.splice(index, 1)
      row.Disabled = false
    },
    resetForm() {
      this.postForm.Name = ''
      this.postForm.Tel = ''
      this.postForm.DeliveryAddress = ''
      this.postForm.DeliveryTime = ''
      this.postForm.Amount = null
      this.postForm.Deposit = null
      this.postForm.Remarks = ''
      this.postForm.Goods = []
      for (let i = 0; i < this.commodities.length; i++) {
        this.commodities[i].Disabled = false
      }

      this.commodities = []
    }
  }
}
</script>
