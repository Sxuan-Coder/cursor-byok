import { createRouter, createWebHashHistory } from "vue-router";
import Home from "@/views/Home.vue";
import ModelConfig from "@/views/ModelConfig.vue";
import Records from "@/views/Records.vue";
import Settings from "@/views/Config.vue";
import Usage from "@/views/Usage.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: Home,
      meta: { title: "Cursor助手", directlyClose: false },
    },
    {
      path: "/records",
      component: Records,
      meta: { title: "记录", directlyClose: false },
    },
    {
      path: "/model-config",
      component: ModelConfig,
      meta: { title: "模型配置", directlyClose: false },
    },
    {
      path: "/usage",
      component: Usage,
      meta: { title: "用量报表", directlyClose: false },
    },
    {
      path: "/settings",
      component: Settings,
      meta: { title: "设置", directlyClose: false },
    },
  ],
});

export default router;