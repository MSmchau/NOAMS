<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>{{ isEdit ? '编辑设备' : '添加设备' }}</template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 700px;"
      >
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入设备名称" />
        </el-form-item>

        <el-form-item label="管理IP" prop="management_ip">
          <el-input v-model="form.management_ip" placeholder="请输入管理IP地址" />
        </el-form-item>

        <el-form-item label="SSH端口">
          <el-input-number v-model="form.ssh_port" :min="1" :max="65535" />
        </el-form-item>

        <el-form-item label="设备类型" prop="device_type">
          <el-select v-model="form.device_type" style="width: 100%;">
            <el-option label="H3C Comware" value="hp_comware" />
            <el-option label="华为 VRP" value="huawei" />
            <el-option label="思科 IOS" value="cisco_ios" />
            <el-option label="锐捷" value="ruijie_os" />
          </el-select>
        </el-form-item>

        <el-form-item label="厂商">
          <el-select v-model="form.vendor" style="width: 100%;">
            <el-option label="H3C" value="h3c" />
            <el-option label="华为" value="huawei" />
            <el-option label="思科" value="cisco" />
            <el-option label="锐捷" value="ruijie" />
          </el-select>
        </el-form-item>

        <el-form-item label="设备角色">
          <el-select v-model="form.role" style="width: 100%;">
            <el-option label="核心交换机" value="core" />
            <el-option label="无线控制器(AC)" value="ac" />
            <el-option label="本体交换机" value="aggregation" />
            <el-option label="接入交换机" value="access" />
            <el-option label="无线AP" value="ap" />
          </el-select>
        </el-form-item>

        <el-form-item label="设备型号">
          <el-input v-model="form.model" placeholder="如 S5560X-30F-EI" />
        </el-form-item>

        <el-form-item label="所属楼栋">
          <el-input v-model="form.building" placeholder="如 宿舍楼1" />
        </el-form-item>

        <el-form-item label="楼层">
          <el-input-number v-model="form.floor" :min="0" />
        </el-form-item>

        <el-form-item label="AP名称" v-if="form.role === 'ap'">
          <el-input v-model="form.ap_name" placeholder="AP在AC上的名称" />
        </el-form-item>

        <el-form-item label="SNMP团体字">
          <el-input v-model="form.snmp_community" placeholder="可选" />
        </el-form-item>

        <el-form-item label="设备描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ isEdit ? '保存修改' : '添加设备' }}
          </el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDevice, createDevice, updateDevice } from '@/api/device'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()

const isEdit = computed(() => !!route.params.id)
const submitting = ref(false)

const form = reactive({
  name: '',
  management_ip: '',
  ssh_port: 22,
  device_type: 'hp_comware',
  vendor: 'h3c',
  role: 'access',
  model: '',
  building: '',
  floor: 0,
  ap_name: '',
  snmp_community: '',
  description: '',
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  management_ip: [
    { required: true, message: '请输入管理IP', trigger: 'blur' },
    { pattern: /^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/, message: '请输入有效的IP地址', trigger: 'blur' },
  ],
  device_type: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
}

async function loadDevice() {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const res = await getDevice(id)
    const d = res.data
    Object.assign(form, {
      name: d.name || '',
      management_ip: d.management_ip || '',
      ssh_port: d.ssh_port || 22,
      device_type: d.device_type || 'hp_comware',
      vendor: d.vendor || 'h3c',
      role: d.role || 'access',
      model: d.model || '',
      building: d.building || '',
      floor: d.floor || 0,
      ap_name: d.ap_name || '',
      snmp_community: '',
      description: d.description || '',
    })
  } catch {
    ElMessage.error('加载设备信息失败')
    router.push('/devices')
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateDevice(Number(route.params.id), { ...form })
      ElMessage.success('设备已更新')
    } else {
      await createDevice({ ...form })
      ElMessage.success('设备已添加')
    }
    router.push('/devices')
  } catch {
    // handled by interceptor
  } finally {
    submitting.value = false
  }
}

function handleCancel() {
  router.push('/devices')
}

onMounted(() => {
  if (isEdit.value) loadDevice()
})
</script>
