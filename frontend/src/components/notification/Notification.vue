<template>

    <q-page class="constrain q-pa-md">
        <div class="row q-col-gutter-lg">
            <div class="col-12">
                <q-list bordered padding >
                    <div v-for="notify in NotifyList" :key="notify._id">
                        <q-item clickable @click="MoveToThePath(notify)" 
                        :class="{'text-red': !notify.isreaded}">
                        <q-item-section top avatar>
                            <q-avatar v-if="notify?.user?.avatar">
                                <img :src="notify?.user?.avatar">
                            </q-avatar>
                            <q-avatar v-else>
                                <img src="https://cdn-icons-png.flaticon.com/512/3237/3237472.png">
                            </q-avatar>
                        </q-item-section> 

                        <q-item-section>
                            <q-item-label>{{ notify.details }}</q-item-label>
                            <q-item-label>{{ notify?.user?.name }}</q-item-label>
                        </q-item-section> 
                    </q-item>
                    <q-separator spaced />
                    </div>
                </q-list>
            </div>
        </div>
    </q-page>
</template>


<script>

import {mapActions, mapGetters, mapState} from 'vuex'


export default {
    name: 'Notification-Component',
    data() {
        return {
            NotifyList: []
        }
    },
    watch: {
        "RealTimeNotify.notifyidData": async function (notify) {
            console.log("noty", notify)
            this.NotifyList.unshift(notify)
        }
    },
    async mounted() {
        var id = this.GetUserData()?.result?._id
        this.NotifyList = await this.GetUnReadedNotifyNum(id)

        setTimeout(() => {
            this.NotifyList.forEach(async el => {
                if (!el.isreaded) {
                    await this.MarkNotificationAsReaded(id)
                    el.isreaded = true
                }
            })
        }, 500)

        // mark notification as readed
    },
    computed: {
        ...mapGetters(['GetUserData']),
        ...mapState(['RealTimeNotify']),
    },
    methods: {
        ...mapActions(['GetUnReadedNotifyNum', 'MarkNotificationAsReaded']),
        MoveToThePath(notify) {
            if (notify?.details.toString().includes("Post")) {
                this.$router.push(`/PostDetails/${notify?.targetid}`)
            } else {
                this.$router.push(`/Profile/${notify?.targetid}`)
            }
        }
    }
}


</script>