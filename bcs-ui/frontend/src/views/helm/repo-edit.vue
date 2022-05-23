<template>
    <div v-bkloading="{ isLoading: loading }">
        <ContentHeader :title="title"></ContentHeader>
        <div class="repo-edit">
            <bk-form :label-width="labelWidth" :model="data" :rules="formRules" ref="formRef">
                <bk-form-item :label="$t('名称')" property="name" error-display-type="normal" required>
                    <bk-input v-model="data.name"></bk-input>
                </bk-form-item>
                <bk-form-item :label="$t('描述')">
                    <bk-input type="textarea" v-model="data.description"></bk-input>
                </bk-form-item>
                <bk-form-item label="Index URL" property="url" error-display-type="normal" required>
                    <bk-input type="textarea" v-model="data.url"></bk-input>
                </bk-form-item>
                <bk-form-item :label="$t('认证')">
                    <bk-radio-group v-model="isAuth">
                        <bk-radio :value="false">{{$t('不认证')}}</bk-radio>
                        <bk-radio :value="true">Basic Auth</bk-radio>
                    </bk-radio-group>
                </bk-form-item>
                <template v-if="isAuth">
                    <bk-form-item :label="$t('用户名')" property="username" error-display-type="normal" required>
                        <bk-input v-model="data.username"></bk-input>
                    </bk-form-item>
                    <bk-form-item :label="$t('密码')" property="password" error-display-type="normal" required>
                        <bk-input v-model="data.password"></bk-input>
                    </bk-form-item>
                </template>
                <bk-form-item>
                    <bk-button theme="primary"
                        :loading="saveLoading"
                        class="bcs-primary-btn"
                        @click="handleSaveData">{{ $t('保存') }}</bk-button>
                    <bk-button class="bcs-primary-btn" @click="handleCancel">{{ $t('取消') }}</bk-button>
                </bk-form-item>
            </bk-form>
        </div>
    </div>
</template>
<script lang="ts">
    import { computed, defineComponent, onMounted, ref } from '@vue/composition-api'
    import useFormLabel from '@/common/use-form-label'
    import $i18n from '@/i18n/i18n-setup'
    import $router from '@/router/index'
    import $store from '@/store'
    import ContentHeader from '@/views/content-header.vue'

    export default defineComponent({
        components: { ContentHeader },
        props: {
            name: {
                type: String,
                default: ''
            }
        },
        setup (props, ctx) {
            const { $bkMessage } = ctx.root

            const formRef = ref<any>()
            const data = ref({
                name: '',
                url: '',
                description: '',
                username: '',
                password: ''
            })
            const isAuth = ref(false)
            const formRules = ref({
                name: [
                    {
                        required: true,
                        message: $i18n.t('必填项'),
                        trigger: 'blur'
                    }
                ],
                url: [
                    {
                        required: true,
                        message: $i18n.t('必填项'),
                        trigger: 'blur'
                    }
                ],
                username: [
                    {
                        required: true,
                        message: $i18n.t('必填项'),
                        trigger: 'blur'
                    }
                ],
                password: [
                    {
                        required: true,
                        message: $i18n.t('必填项'),
                        trigger: 'blur'
                    }
                ]
            })
            const loading = ref(false)
            const saveLoading = ref(false)
            const title = computed(() => {
                return props.name ? $i18n.t('编辑') + props.name : $i18n.t('创建仓库')
            })

            const handleGetDetail = async () => {
                if (!props.name) return
                loading.value = true
                const res = await $store.dispatch('helm/helmRepoDetail', {
                    $repoName: props.name
                })
                data.value = {
                    name: res.name,
                    url: res.configuration?.proxy?.channelList?.[0]?.url,
                    description: res.description,
                    username: res.configuration?.proxy?.channelList?.[0]?.username,
                    password: res.configuration?.proxy?.channelList?.[0]?.password
                }
                if (data.value.password && data.value.username) {
                    isAuth.value = true
                }
                loading.value = false
            }

            const handleSaveData = async () => {
                const valid = await formRef.value?.validate()
                if (!valid) return

                saveLoading.value = true
                let result = false
                const params = isAuth.value ? data.value : {
                    name: data.value.name,
                    url: data.value.url,
                    description: data.value.description
                }
                if (props.name) {
                    result = await $store.dispatch('helm/editHelmRepo', {
                        ...params,
                        $repoName: props.name
                    })
                    result && $bkMessage({
                        theme: 'success',
                        message: $i18n.t('编辑成功')
                    })
                } else {
                    result = await $store.dispatch('helm/createHelmRepo', params)
                    result && $bkMessage({
                        theme: 'success',
                        message: $i18n.t('创建成功')
                    })
                }
                
                result && $router.push({ name: 'helmRepo' })
                saveLoading.value = false
            }

            const handleCancel = () => {
                $router.back()
            }

            const { labelWidth, initFormLabelWidth } = useFormLabel()
            onMounted(() => {
                initFormLabelWidth(formRef.value)
                handleGetDetail()
            })
            return {
                loading,
                saveLoading,
                title,
                isAuth,
                formRef,
                data,
                formRules,
                labelWidth,
                handleCancel,
                handleSaveData
            }
        }
    })
</script>
<style lang="postcss" scoped>
.repo-edit {
    padding: 20px;
    /deep/ .bk-form-control {
        max-width: 640px;
    }
    /deep/ .user-info {
        display: flex;
        max-width: 640px;
    }
}
</style>
