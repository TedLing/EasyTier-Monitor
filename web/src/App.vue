<template>
  <header class="nav-bar">
    <div class="nav-inner">
      <div class="nav-brand">
        <img src="/img/easytier.png" alt="EasyTier" class="nav-logo" />
        <span>EasyTier Monitor</span>
      </div>
      <nav class="nav-links">
        <router-link to="/" class="nav-link" active-class="active" exact-active-class="active">节点监控</router-link>
        <router-link to="/status" class="nav-link" active-class="active">运行状态</router-link>
        <router-link to="/connection" class="nav-link" active-class="active">连接信息</router-link>
      </nav>
    </div>
  </header>
  <router-view />
</template>

<style scoped>
.nav-bar {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-inner {
  max-width: 1500px;
  margin: 0 auto;
  padding: 0 20px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.nav-logo {
  width: 28px;
  height: 28px;
}

.nav-links {
  display: flex;
  gap: 8px;
}

.nav-link {
  padding: 6px 16px;
  border-radius: 6px;
  font-size: 14px;
  color: #666;
  text-decoration: none;
  transition: all 0.2s;
}

.nav-link:hover {
  background: #f0f0f0;
  color: #333;
}

.nav-link.active {
  background: rgba(33, 150, 243, 0.1);
  color: #2196f3;
}
</style>
