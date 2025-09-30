const Layout = () => import('@/layout/index.vue')

export default {
  name: 'Diary',
  path: '/diary',
  component: Layout,
  redirect: '/diary/list',
  meta: {
    title: '日记列表',
    icon: 'mdi:math-log',
    order: 6,
  },
  children: [
    {
      name: 'DiaryList',
      path: 'list',
      component: () => import('./list/index.vue'),
      meta: {
        title: '日记列表',
        icon: 'mdi:book-open-page-variant-outline',
        keepAlive: true,
      },
    }
  ],
}
