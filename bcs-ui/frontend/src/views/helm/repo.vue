<template>
    <div class="helm-repo">
        <div class="helm-repo-header">
            <bk-button class="bcs-primary-btn"
                theme="primary"
                icon="plus"
                @click="handleCreateRepo">{{ $t('创建') }}</bk-button>
            <bk-input right-icon="bk-icon icon-search"
                class="helm-repo-search"
                :placeholder="$t('输入名称搜索')"
                v-model="searchValue"
            ></bk-input>
        </div>
        <bcs-table class="mt20"
            :data="curPageData"
            :pagination="pagination"
            v-bkloading="{ isLoading: loading }"
            @page-change="pageChange"
            @page-limit-change="pageSizeChange">
            <bcs-table-column :label="$t('名称')">
                <template #default="{ row }">
                    <bk-button text @click="handleShowDetail(row)">{{ row.name }}</bk-button>
                </template>
            </bcs-table-column>
            <!-- <bcs-table-column :label="$t('状态')"></bcs-table-column> -->
            <bcs-table-column :label="$t('URL地址')">
                <template #default="{ row }">
                    {{getRepoURL(row)}}
                </template>
            </bcs-table-column>
            <bcs-table-column :label="$t('创建人')" prop="createdBy"></bcs-table-column>
            <bcs-table-column :label="$t('创建时间')" prop="createdDate"></bcs-table-column>
            <bcs-table-column :label="$t('操作')" width="180">
                <template #default="{ row }">
                    <div v-bk-tooltips="{ disabled: row.is_imported, content: $t('非纳管仓库，不允许操作') }">
                        <bk-button text :disabled="!row.is_imported" @click="handleEditRepo(row)">{{$t('编辑')}}</bk-button>
                        <bk-button text
                            :disabled="!row.is_imported"
                            class="ml10"
                            @click="handleRefreshRepo(row)">{{$t('刷新')}}</bk-button>
                        <bk-button text :disabled="!row.is_imported" class="ml10" @click="handleDeleteRepo(row)">{{$t('删除')}}</bk-button>
                    </div>
                </template>
            </bcs-table-column>
        </bcs-table>
        <bk-sideslider
            quick-close
            :is-show.sync="showDetail"
            :title="detail.name"
            :width="640">
            <template #content>
                <RepoDetail :detail="detail"></RepoDetail>
            </template>
        </bk-sideslider>
    </div>
</template>
<script lang="ts">
    import { defineComponent, onMounted, ref } from '@vue/composition-api'
    import useSearch from '@/views/dashboard/common/use-search'
    import usePage from '@/views/dashboard/common/use-page'
    import $store from '@/store'
    import $i18n from '@/i18n/i18n-setup'
    import $router from '@/router/index'
    import RepoDetail from './repo-review.vue'

    export default defineComponent({
        components: { RepoDetail },
        setup (props, ctx) {
            const { $bkInfo, $bkMessage } = ctx.root
            const loading = ref(false)
            const data = ref<any[]>([])
            const keys = ref(['name'])
            const { searchValue, tableDataMatchSearch } = useSearch(data, keys)
            const {
                pagination,
                curPageData,
                pageChange,
                pageSizeChange
            } = usePage(tableDataMatchSearch)

            // 仓库列表
            const handleGetRepoList = async () => {
                loading.value = true
                data.value = await $store.dispatch('helm/helmRepoList')
                loading.value = false
            }
            // 创建仓库
            const handleCreateRepo = () => {
                $router.push({ name: 'helmEdit' })
            }
            // 编辑仓库
            const handleEditRepo = (row) => {
                $router.push({ name: 'helmEdit', params: { name: row.name } })
            }
            // 刷下仓库
            const handleRefreshRepo = async (row) => {
                loading.value = true
                const result = await $store.dispatch('helm/refreshHelmRepo', {
                    $repoName: row.name
                })
                if (result) {
                    $bkMessage({
                        theme: 'success',
                        message: $i18n.t('刷新成功')
                    })
                    await handleGetRepoList()
                }
                loading.value = false
            }
            // 删除仓库
            const handleDeleteRepo = (row) => {
                $bkInfo({
                    type: 'warning',
                    clsName: 'custom-info-confirm',
                    subTitle: row.name,
                    title: $i18n.t('确定删除'),
                    defaultInfo: true,
                    confirmFn: async (vm) => {
                        const result = await $store.dispatch('helm/deleteHelmRepo', {
                            $repoName: row.name
                        })
                        if (result) {
                            $bkMessage({
                                theme: 'success',
                                message: $i18n.t('删除成功')
                            })
                            handleGetRepoList()
                        }
                    }
                })
            }
            const getRepoURL = (row) => {
                return row.configuration?.proxy?.channelList?.[0]?.url || '--'
            }
            // 显示详情
            const showDetail = ref(false)
            const detail = ref({})
            const handleShowDetail = (row) => {
                showDetail.value = true
                detail.value = row
            }

            onMounted(() => {
                handleGetRepoList()
            })
            return {
                loading,
                searchValue,
                pagination,
                curPageData,
                showDetail,
                detail,
                pageChange,
                pageSizeChange,
                handleCreateRepo,
                handleEditRepo,
                handleDeleteRepo,
                handleRefreshRepo,
                getRepoURL,
                handleShowDetail
            }
        }
    })
</script>
<style lang="postcss" scoped>
.helm-repo {
    padding: 20px;
    &-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    &-search {
        width: 400px;
    }
}
</style>
