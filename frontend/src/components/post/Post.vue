<template>
    <div>
        <q-card v-if="!EditPost" class="card-localPost q-mb-md" flat bordered>
            <q-item>
                <q-item-section avatar>
                    <q-avatar>
                        <img v-if="user?.imageUrl" :src="user?.imageUrl"/>
                        <img v-else src="https://cdn-icons-png.flaticon.com/512/1077/1077063.png"/>

                    </q-avatar>
                </q-item-section>

                <q-item-section>
                    <q-item-label class="text-bold">{{user?.name}}</q-item-label>
                    <q-item-label caption>{{getTime()}}</q-item-label>
                </q-item-section>
            </q-item>
            <q-separator />
            <q-img style="cursor: pointer;" @click="GoToDetails" :src="postData.selectedFile" />
            <q-card-section>
                <div class="text-h6">{{postData.title}}</div>
                <div class="text-subtitle1">{{postData.message}}</div>
                <q-separator />
                <div class="text-subtitle4" 
                v-for="(comment, index) in postData.comments" :key="index">
                    {{comment}}
                </div>
                <q-btn v-if="!UserLike" @click="Like" flat round color="red" icon="eva-heart-outline">
                    {{LikesCount()}}
                </q-btn>

                <q-btn v-else @click="Like" flat round color="red" icon="eva-heart">
                    {{LikesCount()}}
                </q-btn>
            </q-card-section>

            <q-input outlined v-model="form.text" label = "add comment..">
                <q-btn v-if="form.text !== ''" @click="AddComment" flat round color="primary" icon="eva-plus-square" />
            </q-input>
        </q-card>
        <div v-else class="q-pa-md items-start q-gutter-md">
            <q-card class="my-card col-12">
                <q-card-section>
                    <div class="text-h6">Edit Post</div>
                    <q-input dense v-model="postData.title" autofocus placeholder="Post Title"/>
                    <div>
                        <q-input v-model="postData.message" placeholder="What's on your mind?" type="textarea"/>
                    </div>
                    <div class="q-pa-md">
                        <q-file v-model="file" label="Pick Image" filled/>
                    </div>
                    <div>
                        <q-img :src="postData.selectedFile" spinner-color="red" style="height: 140px; max-width: 150px;"/>
                    </div>

                    <q-btn flat label="Update" v-close-popup @click="FireUpdate"/>
                </q-card-section>
            </q-card>
        </div>
    </div>
</template>

<script>
import {mapActions, mapGetters} from 'vuex'
import moment from 'moment'
export default {
    name: 'PostComponent',
    props: ['post', 'EditPost'],
    data() {
        return {
            user: {},
            form: {text: ''},
            file: null,
            UserLike: false,
            postData: {}
        }
    },
    watch: {
        file() {
            this.ConvertToBase64()
        }
    }, 
    methods: {
        ...mapActions(['GetUserByID', 'LikePostByUser', 'commentPost', 'updatePost']),

        GoToDetails() {
            this.$router.push({path: `/PostDetails/${this.postData?._id}`})
        },

        async FireUpdate() {
            const PostData = {
                id: this.postData?._id,
                title: this.postData?.title,
                message: this.postData?.message,
                selectedFile: this.postData?.selectedFile,
            }
            const res = await this.updatePost(PostData)
            if (res) {
                this.$emit('changeEdit')
            }
        },

        getTime() {
            return moment(this.postData?.createdAt).fromNow()
        },

        Like() {
            this.LikePostByUser(this.postData?._id)
            const userData = this.GetUserData()
            if (userData?.result?._id) {
                const uid = userData.result._id
                if (this.UserLike) {
                    this.postData.likes = this.postData.likes.filter(id => id != uid)
                } else {
                    this.postData.likes.push(uid)
                }
                this.UserLike = !this.UserLike
            }
        },
        LikesCount() {
            if (this.postData?.likes?.length > 0) {
                return String(this.postData?.likes?.length)
            }
        },

        AddComment() {
            this.postData.comments.push(this.form.text)
            // store
            this.commentPost({Value: this.form.text, id: this.postData?._id})
            this.form.text = ''
        },

        ConvertToBase64() {
            var reader = new FileReader()
            reader.readAsDataURL(this.file)
            reader.onloadend = () => {
                this.postData.selectedFile = reader.result
            }
        },
    },
    computed: {
        ...mapGetters(['GetUserData'])
    },
    async mounted() {
        if (!this.post) return
        this.postData = JSON.parse(JSON.stringify(this.post))

        // 如果 creator 已经是用户对象，直接使用
        if (this.postData?.creator?.name) {
            this.user = this.postData.creator
        } else if (this.postData?.creator && typeof this.postData.creator === 'string' && this.postData.creator.length > 10) {
            // 只有 creator 是有效字符串 ID 时才调用 API
            const res = await this.GetUserByID(this.postData.creator)
            this.user = res?.user || {}
        }

        // get if user liked the postData or not
        const userData = this.GetUserData()
        if (userData?.result?._id && this.postData?.likes) {
            const uid = userData.result._id
            this.UserLike = this.postData.likes.some((like) => like == uid)
        }
    }
}
</script>