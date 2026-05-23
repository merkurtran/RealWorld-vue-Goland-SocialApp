import * as api from '../api/index.js'

import { jwtDecode } from 'jwt-decode'

const Auth = {
  state: {
    authData: null
  },
  getters: {
    GetUserData: (state) => () => {
      return state.authData
    }
  },
  mutations: {
    Auth(state, payload) {
        localStorage.setItem('profile', JSON.stringify({...payload}))
        state.authData = payload
    },
    SetData(state) {
        const user = JSON.parse(localStorage.getItem('profile'))
        if (!user) return

        const token = user?.token
        if (token) {
            try {
                const decodedToken = jwtDecode(token)
                // 检查 token 是否过期
                if (decodedToken.exp && decodedToken.exp * 1000 < new Date().getTime()) {
                    this.commit('Logout')
                    return
                }
            } catch (e) {
                // token 无效
                this.commit('Logout')
                return
            }
        }
        state.authData = user
    },
    Logout(state) {
        localStorage.clear()
        state.authData = null
    }
  },
  actions: {
    async signin (context, formData) {
        try {
            const {data} = await api.signIn(formData)
            context.commit('Auth', data)
            // context.commit('SetData')
            return data
        } catch (error) {
            console.log(error)
            return error
        }
    },
    async signup (context, formData) {
        try {
            const {data} = await api.signUp(formData)
            context.commit('Auth', data)
            context.commit('SetData')
            return data
        } catch (error) {
            console.log(error)
            return error
        }
    },
    async logout (context) {
        try {
            context.commit('Logout')
        } catch (error) {
            console.log(error)
            return error
        }
    }
  },
}


export default Auth