<template>
    <bcs-navigation navigation-type="top-bottom" :need-menu="false">
        <template slot="side-header">
            <span class="title-icon"><img src="@/images/bcs.svg" class="all-icon"></span>
            <span class="title-desc bcs-title-desc" @click="handleGoHome">{{ $INTERNAL ? $t('TKEx-IEG 容器平台') : $t('蓝鲸容器管理平台') }}</span>
        </template>
        <template #header>
            <div class="bcs-navigation-header">
                <div class="nav-left">
                    <bcs-select ref="projectSelectRef" class="header-select" :clearable="false" searchable
                        :value="curProjectCode"
                        v-show="$route.name !== 'projectManage'"
                        @change="handleProjectChange">
                        <bcs-option v-for="option in onlineProjectList"
                            :key="option.project_code"
                            :id="option.project_code"
                            :name="option.project_name">
                        </bcs-option>
                        <template #extension v-if="!$INTERNAL">
                            <div class="extension-item" @click="handleCreateProject"><i class="bk-icon icon-plus-circle mr5"></i>{{$t('新建项目')}}</div>
                            <div class="extension-item" @click="handleGotoProjectManage"><i class="bcs-icon bcs-icon-apps mr5"></i>{{$t('项目管理')}}</div>
                        </template>
                    </bcs-select>
                </div>
                <div class="nav-right">
                    <div class="header-help" @click="handleGotoHelp">
                        <i class="bcs-icon bcs-icon-help-2"></i>
                    </div>
                    <bcs-popover theme="light navigation-message" :arrow="false" offset="0, 10" placement="bottom-start" :tippy-options="{ 'hideOnClick': false }">
                        <div class="header-user">
                            {{user.username}}
                            <i class="bk-icon icon-down-shape"></i>
                        </div>
                        <template slot="content">
                            <ul class="bcs-navigation-admin">
                                <li class="nav-item" @click="handleGotoProjectManage">{{ $t('项目管理') }}</li>
                                <li class="nav-item" @click="handleLogout">{{ $t('退出') }}</li>
                            </ul>
                        </template>
                    </bcs-popover>
                </div>
            </div>
        </template>
        <template #default>
            <slot></slot>
        </template>
    </bcs-navigation>
</template>
<script>
    export default {
        name: "Navigation",
        data () {
            return {
            }
        },
        computed: {
            user () {
                return this.$store.state.user
            },
            onlineProjectList () {
                return this.$store.state.sideMenu.onlineProjectList
            },
            curProjectCode () {
                return this.$store.state.curProjectCode
            },
            curCluster () {
                const cluster = this.$store.state.cluster.curCluster
                return cluster && Object.keys(cluster).length ? cluster : null
            },
            curProject () {
                return this.$store.state.curProject
            }
        },
        methods: {
            async handleProjectChange (code) {
                // 解决组件初始化时触发change事件问题
                if (code === this.curProjectCode) return

                const item = this.onlineProjectList.find(item => item.project_code === code)
                if (item?.kind !== this.curProject.kind) {
                    // 切换不同项目时刷新界面
                    const route = this.$router.resolve({
                        name: 'clusterMain',
                        params: {
                            projectCode: code,
                            // eslint-disable-next-line camelcase
                            projectId: item?.project_id
                        }
                    })
                    location.href = route.href
                } else {
                    this.$router.push({
                        name: 'clusterMain',
                        params: {
                            projectCode: code,
                            // eslint-disable-next-line camelcase
                            projectId: item?.project_id
                        }
                    })
                }
            },
            handleGotoProjectManage () {
                this.$refs.projectSelectRef && this.$refs.projectSelectRef.close()
                if (window.REGION === 'ieod') {
                    window.open(`${window.DEVOPS_HOST}/console/pm`)
                } else {
                    if (this.$route.name === 'projectManage') return
                    this.$router.push({
                        name: 'projectManage'
                    })
                }
            },
            handleCreateProject () {
                this.$refs.projectSelectRef && this.$refs.projectSelectRef.close()
                this.$emit('create-project')
            },
            handleGotoHelp () {
                window.open(window.BCS_CONFIG?.doc?.help)
            },
            // 跳转首页
            handleGoHome () {
                if (this.$route.name !== 'clusterMain' && !this.curCluster) {
                    // 全部集群首页
                    this.$router.push({ name: 'clusterMain' })
                } else if (this.$route.name !== 'clusterOverview' && this.curCluster) {
                    // 单集群首页
                    this.$router.replace({ name: 'clusterOverview' })
                }
            },
            // 注销
            handleLogout () {
                window.location.href = `${LOGIN_FULL}?c_url=${window.location}`
            }
        }
    }
</script>
<style lang="postcss" scoped>
/deep/ .bk-navigation-wrapper .container-content {
    padding: 0;
}
/deep/ .bk-select .bk-tooltip.bk-select-dropdown {
    background: transparent;
}
/deep/ .bcs-title-desc {
    cursor: pointer;
}
.all-icon {
    width: 28px;
    height: 28px;
}
.bcs-navigation-admin {
    display:flex;
    flex-direction:column;
    background:#FFFFFF;
    border:1px solid #E2E2E2;
    margin:0;
    color:#63656E;
    padding: 6px 0;
}
.nav-item {
    flex:0 0 32px;
    display:flex;
    align-items:center;
    padding:0 20px;
    list-style:none;
    &:hover {
        color:#3A84FF;
        cursor:pointer;
        background-color:#F0F1F5;
    }
}
.extension-item {
    margin: 0 -16px;
    padding: 0 16px;
    &:hover {
        cursor: pointer;
        background-color: #f0f1f5;
    }
}
/deep/ .create-input {
    width: 90%;
}
.bcs-navigation-header {
    flex:1;
    height:100%;
    display:flex;
    align-items:center;
    justify-content: space-between;
    font-size:14px;
    .nav-left {
        flex: 1;
        display:flex;
        padding:0;
        margin:0;
        .header-select {
            width:240px;
            margin-right:34px;
            border:none;
            background:#252F43;
            color:#D3D9E4;
            box-shadow:none;
        }
        .header-nav-item {
            list-style:none;
            height:50px;
            display:flex;
            align-items:center;
            margin-right:40px;
            color:#96A2B9;
            min-width:56px;
            &:hover {
                cursor:pointer;
                color:#D3D9E4;
            }
            &.active {
                color: #fff;
            }
        }
    }
    .nav-right {
        display: flex;
        align-items: center;
        .header-help {
            color:#768197;
            font-size:16px;
            position:relative;
            height:32px;
            width:32px;
            display:flex;
            align-items:center;
            justify-content:center;
            margin-right:8px;
            &:hover {
                background:linear-gradient(270deg,rgba(37,48,71,1) 0%,rgba(38,50,71,1) 100%);
                border-radius:100%;
                cursor:pointer;
                color:#D3D9E4;
            }
        }
        /deep/ .header-user {
            height:100%;
            display:flex;
            align-items:center;
            justify-content:center;
            color:#96A2B9;
            margin-left:8px;
            .bk-icon {
                margin-left:5px;
                font-size:12px;
            }
            &:hover {
                cursor:pointer;
                color:#D3D9E4;
            }
        }
    }
}
</style>
