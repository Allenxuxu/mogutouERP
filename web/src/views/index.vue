<template>
  <div class="dashboard-container">
    <el-carousel arrow="never">
      <div v-if="device === 'desktop'">
        <el-carousel-item style="text-align: center">
          <img src="../assets/log/post.png" width="1000px" height="100%">
        </el-carousel-item>
      </div>
      <div v-else>
        <el-carousel-item style="text-align: center">
          <img src="../assets/log/log.png">
        </el-carousel-item>
      </div>
    </el-carousel>

    <div style="text-align: center" class="dashboard-text">欢迎您</div>
    <div style="text-align: center" class="dashboard-text">{{ name }}</div>
    <br>
    <el-row style="text-align: center">
      <router-link class="inlineBlock" to="/stock/index">
        <el-button type="primary" icon="el-icon-goods" round>库存查询</el-button>
      </router-link>
      <router-link class="inlineBlock" to="/order/sales/create">
        <el-button type="success" icon="el-icon-circle-plus-outline" round>客户订单</el-button>
      </router-link>
    </el-row>
    <br>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Dashboard',
  computed: {
    ...mapGetters(['name', 'device'])
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
