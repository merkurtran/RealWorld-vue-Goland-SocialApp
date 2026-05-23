<template>
  <q-page class="constrain q-pa-md">
    <div class="row q-col-gutter-lg">
        <div class="col-3"></div>
        <div class="col-6" v-if="post">
            <div class="q-pa-md q-gutter-sm" v-if="!EditPost">
                <q-btn v-if="IsSameUser" color="primary" icon="eva-edit" @click="EditPost = !EditPost" label="Edit Post"/>
            </div>
            <Post :post="post" :EditPost="EditPost" @changeEdit="EditPost = !EditPost"/>
        </div>
        <div class="col-3"></div>
    </div>
  </q-page>
</template>

<script>
import Post from './Post.vue';
import { mapActions } from 'vuex';

export default {
  name: 'PostDetails',
  data() {
    return {
        EditPost: false,
        post: null,
        IsSameUser: false
    }
  },
  methods: {
    ...mapActions(['getPost'])
  },
  async mounted() {
    const id = this.$route.params.id
    const {post} = await this.getPost(id)
    this.post = post
    const logedInUser = JSON.parse(localStorage.getItem('profile'))?.result?._id
    if (logedInUser === this.post?.creator) {
        this.IsSameUser = true
    }
  },
  commponents: {
    Post
  }
}
</script>