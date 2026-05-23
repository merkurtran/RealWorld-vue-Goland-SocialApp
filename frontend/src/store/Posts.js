import * as api from "../api/index.js";


const Posts = {
    state: {isLodaing: true, post: [], posts: [], SearchResult: []},
    getters: {
        GetPost: (state) => () => {
            return {...state.post}
        },
        GetAllPost: (state) => () => {
            return {...state.posts}
        },
        GetSearchData: (state) => () => {
            return {...state.SearchResult}
        }
    },
    mutations: {
        Post(state, payload) {
            state.post = payload
        },
        Posts(state, payload) {
            state.posts = payload
        },
        Search(state, payload) {
            state.SearchResult = payload
        }
    },
    actions: {
        async getPost(context, id){
            try {
                let {data} = await api.fetchPost(id)
                context.commit('Post', data)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async getPosts(context, page) {
            try {
                let user = JSON.parse(localStorage.getItem('profile'))
                let userId = user?.result?._id
                if (!userId) {
                    throw new Error('User not authenticated')
                }
                const {data} = await api.fetchPosts(page, userId)
                context.commit('Posts', data)
                return data
            } catch (error) {
                console.log(error)
                throw error
            }
        },
        // this for user & posts
        async getPostsUsersBySearch(context, searchQuery) {
            try {
                const {data} = await api.fetchPostsUsersbySearch(searchQuery)
                context.commit('Search', data)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async createPost(context, post) {
            try {
                const {data} = await api.createPost(post)
                context.commit('Post', data)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async updatePost(context, Data) {
            try {
                let user = JSON.parse(localStorage.getItem('profile'))
                let userId = user?.result?._id

                const PostData = {
                    "title": Data.title,
                    "message": Data.message,
                    "creator": userId,
                    "selectedFile": Data.selectedFile,
                }
                const post = await api.updatePost(Data.id, PostData)
                context.commit('Post', post)
                return post
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async LikePostByUser(context, id) {
            try {
                // const user = JSON.parse(localStorage.getItem('profile'))
                const {data} = await api.likePost(id)
                context.commit('Post', data)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async commentPost(context, form) {
            try {
                const {data} = await api.comment(form.Value, form.id)
                context.commit('Post', data)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async deletePost(context, id) {
            try {
                await api.deletePost(id)
            } catch (error) {
                console.log(error)
                return error
            }
        }
    }
}

export default Posts