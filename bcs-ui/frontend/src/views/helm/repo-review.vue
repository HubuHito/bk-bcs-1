<template>
    <div class="repo-review">
        <bk-form :label-width="100">
            <bk-form-item :label="$t('名称')">{{ detail.name }}</bk-form-item>
            <bk-form-item :label="$t('地址')">{{ firstChannelList.url || '--' }}</bk-form-item>
            <bk-form-item :label="$t('描述')">{{ detail.description || '--' }}</bk-form-item>
            <template v-if="isAuth">
                <bk-form-item :label="$t('用户名')">{{ firstChannelList.username }}</bk-form-item>
                <bk-form-item :label="$t('密码')">{{ firstChannelList.password }}</bk-form-item>
            </template>
        </bk-form>
    </div>
</template>
<script lang="ts">
    import { computed, defineComponent } from '@vue/composition-api'

    export default defineComponent({
        props: {
            detail: {
                type: Object,
                default: () => ({})
            }
        },
        setup (props, ctx) {
            const { detail } = props
            const firstChannelList = computed(() => {
                return detail.configuration?.proxy?.channelList?.[0] || {}
            })
            const isAuth = computed(() => {
                return firstChannelList.value.password && firstChannelList.value.username
            })
            return {
                isAuth,
                firstChannelList
            }
        }
    })
</script>
<style lang="postcss" scoped>
.repo-review {
    padding: 20px;
}
</style>
