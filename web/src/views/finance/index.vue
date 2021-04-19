<template>
  <div>
    <el-select v-model="value" class="select" placeholder="请选择年份" @change="handleChange">
      <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value"/>
    </el-select>
    <div class="chart-container">
      <finance-chart :chart-data="chartData"/>
    </div>
    <div class="chart-container">
      <order-chart :chart-data="chartData"/>
    </div>
  </div>
</template>

<script>
import FinanceChart from './Charts/financeCharts'
import OrderChart from './Charts/orderQuantity'

import { getFinance } from '@/api/finance'

export default {
  name: 'DashboardAdmin',
  components: {
    FinanceChart,
    OrderChart
  },
  data() {
    const now = new Date()
    const year = now.getFullYear()
    return {
      options: [{
        value: year,
        label: year
      }, {
        value: year - 1,
        label: year - 1
      }, {
        value: year - 2,
        label: year - 2
      }],
      value: year,
      chartData: {}
    }
  },
  created() {
    this.handleChange(this.value)
  },
  methods: {
    handleChange(v) {
      getFinance(v).then(response => {
        console.log(v, response)
        this.chartData = response
      })
    }
  }
}
</script>

<style scoped>
.select {
  margin-top: 10px;
  margin-left: 30px;
  margin-bottom: 10px;
}
.chart-container {
  position: relative;
  width: 100%;
  height: calc(100vh - 65px);
}
</style>
