const RealTimeNotify = {
  state:{
    ws: null,
    // userid: '',
    notifyideslistNumber: 0,
    notifyidData: null,
  },
  getters: {
    // GetTheuserid: (state) => () => {
    //     return state.userid
    // },
    Getnotifyideslist: (state) => () => {
        return state.notifyideslistNumber
    },
  },
  mutations: {
    // SET_USER_ID(state) {
    //     if(JSON.parse(localStorage.getItem('profile'))){
    //         state.userid = JSON.parse(localStorage.getItem('profile')).result._id
    //     }
    // },
    SET_WS(state, ws){
        state.ws = ws
    },
    ADD_NOTIFICATION(state, notify){
        state.notifyideslistNumber = state.notifyideslistNumber + 1
        state.notifyidData = notify
    }
  },
  actions: {
    async connectToNotify(context){
        if(JSON.parse(localStorage.getItem('profile')) && context.state.ws === null) {
            const Userid = JSON.parse(localStorage.getItem('profile')).result._id
            const uri = process.env.VUE_APP_RealTineNotificationUrl
            const ws = new WebSocket(`${uri}${Userid}`)

            ws.onopen = () => {
                console.log('连接成功')
                context.commit('SET_WS', ws)
            }

            ws.onmessage = (event) => {
                console.log('接收到消息', event.data)
                const notify = JSON.parse(event.data)
                context.commit('ADD_NOTIFICATION', notify)
            }

        }
    },
    async StopConnectionToNotify(context) {
        try {
           context.state.ws.close()
           context.commit('SET_WS', null) 
        } catch (error) {
            console.log(error)
        }
    }
  }
}

export default RealTimeNotify;