<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NPopconfirm, NTag, NTabs, NDatePicker ,NRadio, NRadioGroup,NTabPane} from 'naive-ui'
import { useRouter } from 'vue-router'
import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import UploadMore from '@/components/UploadMore.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: '日记列表' })

function updateOrDeleteDiary(ids) {
  extraParams.value.is_delete
    ? api.deleteDiary(ids)
    : api.softDeleteDiary(JSON.parse(ids), true)
}

const {
  modalVisible,
  handleDelete,
  handleAdd,
  modalForm,
  modalFormRef,
  handleView,
} = useCRUD({
  name: '日记',
  doDelete: updateOrDeleteDiary,
  refresh: () => $table.value?.handleSearch(),
  initForm:{
    content: '',
    imgs: [],
    add_time: Date.now(), // 默认当前时间
    status: 1, // 默认公开
  }
})


const router = useRouter()
//时间范围
const range = ref(null)
const btnLoading = ref(false)

const rules = {
  content: {
    required: true,
    message: '内容不能为空',
  },
  add_time: {
    required: true,
    message: '日期不能为空',
  },
    status: {
    required: true,
    message: '发布形式不能为空',
  },
}

//页面初始化加载
onMounted(() => {
  handleChangeTab('all') // 默认查看全部
})

// 保存
async function handleSave() {
  modalFormRef.value?.validate(async (err) => {
    if (!err) {
      btnLoading.value = true
      console.log(modalForm.value)
      try {
        await api.saveOrUpdateDiary(modalForm.value)
        modalVisible.value = false
        $message.success('操作成功!')
        await router.replace({ path: '/diary/list', query: { needRefresh: true } })
      }
      catch (err) {
        console.error(err)
      }
      finally {
        btnLoading.value = false
      }
    }
  })
}

const extraParams = ref({
  is_delete: null, // 未删除 | 回收站
  status: null, // null-all, 1-公开, 2-私密, 3-草稿
})


function handleTimeRangeChange(value) {
  if (value && value.length === 2) {
    queryItems.value.add_time_start = value[0] // 毫秒时间戳
    queryItems.value.add_time_end = value[1]   // 毫秒时间戳
  } else {
    queryItems.value.add_time_start = 0
    queryItems.value.add_time_end = 0
  }
}

const $table = ref(null)
const queryItems = ref({
  content: '',
  add_time_start: 0,
  add_time_end: 0,
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  { title: '内容', key: 'content', width: 150, align: 'center', ellipsis: { tooltip: true } },
  { title: '发布状态', key: 'status', width: 20, align: 'center', ellipsis: { tooltip: true } },
  {
    title: '发布时间',
    key: 'add_time',
    align: 'center',
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.add_time),
          icon: () => h('i', { class: 'i-mdi:update' }),
        },
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 40,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            quaternary: true,
            type: 'info',
            onClick: () => handleView(row),
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-ic:outline-remove-red-eye' }),
          },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  quaternary: true,
                  type: 'error',
                  style: 'margin-left: 15px;',
                },
                {
                  default: () => '删除',
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
                },
              ),
            default: () => h('div', {}, '确定删除该日记吗?'),
          },
        ),
      ]
    },
  },
]


function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_delete = 0
      extraParams.value.status = null
      break
    case 'public':
      extraParams.value.is_delete = 0
      extraParams.value.status = 1
      break
    case 'secret':
      extraParams.value.is_delete = 0
      extraParams.value.status = 2
      break
    case 'delete':
      extraParams.value.is_delete = 1
      extraParams.value.status = null
      break
  }
  $table.value?.handleSearch()
}
</script>

<template>
  <CommonPage title="日记列表">
    <template #action>
      <NButton type="primary" @click="handleAdd()">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        新增日记
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:playlist-remove" />
        </template>
        批量删除
      </NButton>

    </template>

    <NTabs type="line" animated @update:value="handleChangeTab">
      <template #prefix>
        状态
      </template>
      <NTabPane name="all" tab="全部" />
      <NTabPane name="public" tab="公开" />
      <NTabPane name="secret" tab="私密" />
      <NTabPane name="delete" tab="回收站" />
    </NTabs>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getDiaryList"
    >
      <template #queryBar>
        <QueryItem label="内容" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.content"
            clearable
            type="text"
            placeholder="请输入内容"
          />
        </QueryItem>
        <QueryItem label="时间范围" :label-width="100" :content-width="300">
            <n-date-picker v-model:value="range" type="daterange" clearable 
                 @update:value="handleTimeRangeChange"/>
        </QueryItem>

      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
      title="发布日记"
      :loading="btnLoading"
      show-footer
      @save="handleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="100"
        :model="modalForm"
        :rules="rules"
      >
        <NFormItem label="内容" path="content">
          <n-input
            v-model:value="modalForm.content"
            maxlength="200"
            placeholder=""
            type="textarea"
            show-count 
          />
        </NFormItem> 

        <NFormItem label="图片" path="imgs">
          <UploadMore
            v-model:preview="modalForm.imgs"
          />
        </NFormItem>
        <NFormItem label="发布日期" path="add_time">
          <n-date-picker v-model:value="modalForm.add_time" type="datetime" clearable />
        </NFormItem>
          
        <NFormItem label="发布形式" path="status">
          <NRadioGroup v-model:value="modalForm.status" name="radiogroup">
            <NSpace>
              <NRadio :value="1">
                公开
              </NRadio>
              <NRadio :value="2">
                私密
              </NRadio>
            </NSpace>
          </NRadioGroup>
          
        </NFormItem>

      </NForm>
    </CrudModal>


  </CommonPage>
</template>
