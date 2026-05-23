<template>
<div>
    <q-list bordered v-if="UserData && UserData.length > 0">
        <q-item-label header class="text-bold">推荐关注</q-item-label>
        <q-item class="RightBar" v-for="user in UserData" :key="user._id" :to="`/Profile/${user?._id}`">
            <q-item-section avatar>
                <q-avatar>
                    <img v-if="user.imageUrl" :src="user.imageUrl">
                    <img v-else src="https://cdn-icons-png.flaticon.com/512/1077/1077063.png">
                </q-avatar>
            </q-item-section>
            <q-item-section>
                <q-item-label class="text-bold">{{ user?.name }}</q-item-label>
                <q-item-label caption>{{ user?.bio }}</q-item-label>
            </q-item-section>
        </q-item>
    </q-list>
    <q-card v-else class="q-pa-md text-center text-grey">
        <q-card-section>暂无推荐用户</q-card-section>
    </q-card>
</div>

</template>

<script>
import {mapActions, mapGetters} from 'vuex'

export default {
  name: 'RightBar',
  data() {
    return {
      UserData: []
    }
  },
  computed: {
    ...mapGetters(['GetUserData'])
  },
  methods: {
    ...mapActions(['GetTheUserSug'])
  },
  async mounted() {
    let lgoedinUser = this.GetUserData()?.result
    if (lgoedinUser) {
        const {users} = await this.GetTheUserSug(lgoedinUser?._id)
        console.log("RightBar users:", users)

        this.UserData = users
    }
  }
}

</script>

<style lang="sass" scoped>
.RightBar
 cursor: pointer
</style>