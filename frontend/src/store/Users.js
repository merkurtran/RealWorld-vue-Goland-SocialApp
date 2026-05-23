import * as api from '@/api/index.js'

const Users = {
    state: {User: null},
    getters: {
        GetUser: (state) => () => {
            return state.User
        },
        // TODO GetUserFollowersFollowing
        GetUserFollowersFollowing: async () => {
            const raw = localStorage.getItem('profile')
            if (!raw) return []
            const userd = JSON.parse(raw)
            if (!userd?.result) return []
            var followers = userd.result.followers || []
            var following = userd.result.following || []
            const combinedArray = [...followers, ...following]
            const uniqueArray = Array.from(new Set(combinedArray)).filter(id => id && id.length > 5)
            var userdata = []
            for (const uid of uniqueArray) {
                try {
                    const {data} = await api.fetchUserProfile(uid)
                    var user = {"_id": data.user._id, 
                        "name": data.user.name, 
                        "imageUrl": data.user.imageUrl} 
                    userdata.push(user)
                } catch (error) {
                    console.log(`跳过不存在的用户: ${uid}`)
                }
            }
            return userdata
        }
    },
    mutations: {
        UserData(state, payload) {
            state.User = payload?.data
        }
    },
    actions: {
        // getuserbyid
        async GetUserByID(context, id) {
            try{
                const {data} = await api.fetchUserProfile(id)
                context.commit('UserData', data.user)

                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        // update user data
        async UpdateUserDate(context, userData) {
            try{
                const {data} = await api.UpdateUser(userData)
                context.commit('UserData', data.user)

                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        // following user
        async FollowUser(context, ProfileID) {
            try{
                const {data} = await api.following(ProfileID)

                return data
            } catch (error) {
                console.log(error)
                return error
            }
        },
        async GetTheUserSug(context, id) {
            try {
                const {data} = await api.getSugUser(id)
                return data
            } catch (error) {
                console.log(error)
                return error
            }
        }
    }
}

export default Users