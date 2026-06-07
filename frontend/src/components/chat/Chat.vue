<template>
    <q-page class="constrain q-pa-md">
        <div class="row q-col-gutter-lg">
            <div class="col-12 chat-container">
                <div class="user-list">
                    <div class="q-pa-md">
                        <q-toolbar class="bg-primary text-white shadow-1">
                            <q-toolbar-title>Following & Followers</q-toolbar-title>
                        </q-toolbar>

                        <q-list bordered>
                            <q-item @click="selectUser(contact)" 
                            v-for="contact in contacts" 
                            :key="contact._id" class="q-my-sm" clickable v-ripple>
                            <q-item-section avatar>
                                <q-avatar v-if="!contact.imageUrl" color="primary" text-color="white">
                                    {{ contact.name[0] }}
                                </q-avatar>
                                <q-avatar v-else>
                                    <img :src="contact?.imageUrl"/>
                                </q-avatar>

                            </q-item-section>
                            <q-item-section>

                                <q-item-label>{{ contact.name }}</q-item-label>
                            </q-item-section>

                            <q-item-section side v-if="contact.isOnline">
                                <q-badge color="positive" rounded/>
                            </q-item-section>

                            <q-item-section side v-if="contact.unReadedmessage && contact.unReadedmessage > 0">
                                <q-badge color="negative" rounded :label="contact?.unReadedmessage"/>
                            </q-item-section>
                        
                        
                        </q-item>

                        </q-list>
                    </div>
                </div>
                <!-- chat box -->
                 <div class="chat-messages" v-if="selectedUser != null" style="background: white;">
                    <div class="q-pa-md row justify-center" 
                    style="overflow-y: auto; max-width: 400px; max-height: 500px;" ref="messageContainer" @scroll="handlescroll">
                    <div v-for="msg in messageBetweenUsers" :key="msg._id" style="width: 100%;">
                        <q-chat-message :name="msg.sender === MainUserData._id ? MainUserData.name : selectedUser.name" 
                        :avatar="msg.sender === MainUserData._id ? MainUserData.imageUrl : selectedUser.imageUrl" 
                        :text="[msg.content]" :sent="msg.sender === MainUserData._id ? true : false" />
                    </div>
                 </div>
                 <q-separator spaced />
                 <q-input outlined v-model="messageToSend.text" @keyup.enter="Sendmessage" label="write message..">
                    <q-btn v-if="messageToSend.text != ''" @click="Sendmessage" flat round color="primary" icon="eva-arrow-right"/>
                 </q-input>
            </div>
        </div>
        </div>
    </q-page>

</template>


<script>
import { mapGetters, mapActions, mapState} from 'vuex';

export default {
  name: 'ChatComponent',
  data() {
    return {
      messageToSend: {text: ''},
      contacts: [],
      messageBetweenUsers: [],
      messagelistnum: 0,
      selectedUser: null,
      MainUserData: {},
      uniqueOnlineUsers: {}
    }
  },
  computed: {
    ...mapGetters(['GetUserFollowersFollowing', 'GetUserData']),
    ...mapState(['RealTimeChat'])
  },
  watch: {
    "RealTimeChat.onlineFriends": function (online) {
        const onlineFriendsArray = Object.values(online)
        this.uniqueOnlineUsers = Array.from(new Set(onlineFriendsArray))
        this.updateOnlineList()
    },
    "RealTimeChat.privateMessages": function (messages) {
        // Handle private messages
        if (this.contacts.length > 0) {
            this.contacts.forEach((contact) => {
                if (contact._id === messages.sender) {
                    contact.unReadedmessage ++
                }
            })
            if (this.selectedUser && this.selectedUser?._id === messages.sender) {
                this.messageBetweenUsers.push(messages)
                setTimeout(() => {
                    this.scrollDownFunction()
                }, 100)
            }
        }
    }
  },
  async mounted() {
    this.MainUserData = this.GetUserData()?.result
    await this.GetUsList()

    this.uniqueOnlineUsers = Array.from(new Set(Object.values(this.RealTimeChat.onlineFriends)))

    this.updateOnlineList()
    const saved = localStorage.getItem('chat_selected_user')
    if (saved && this.contacts.length > 0) {
      const savedUser = JSON.parse(saved)
      const found = this.contacts.find(c => c._id === savedUser._id)
      if (found) {
        await this.selectUser(found)
      }
    }
  },
  methods: {
    ...mapActions(['GetUnreadedMessageNum', 'GetMsgsBetweebTwoUsers', 'SendMessage', 'MarkMsgsAsReaded', 'SendPrivateMessage']),
    updateOnlineList() {
        this.contacts.forEach((contact) => {
            if(this.uniqueOnlineUsers.includes(contact._id)) {
                contact.isOnline = true 
            } else {
                contact.isOnline = false
            }
        })
    },
    handlescroll() {
        const container = this.$refs.messageContainer
        if (container.scrollTop === 0) {
            this.GetOldestMessagesBetweenUsers()
        }
    },
    async GetOldestMessagesBetweenUsers() {
        this.messagelistnum = this.messagelistnum + 1
        var firstuid = this.MainUserData?._id
        var seconduid = this.selectedUser?._id
        var from = this.messagelistnum
        var ndata = {from, firstuid, seconduid}
        var res = await this.GetMsgsBetweebTwoUsers(ndata)
        if (!res) return
        var {msgs} = res
        if (!msgs) return
        this.messageBetweenUsers.unshift(...msgs)
    },
    scrollDownFunction() {
        const container = this.$refs.messageContainer
        container.scrollTop = container.scrollHeight
    },
    async CallMarkMsgAsReaded(user) {
        var mainuid = this.MainUserData?._id
        var otheruid = user?._id
        var GetunReadedmessage = 0

        this.contacts.forEach(c => {
            if (String(otheruid) == String(c?._id)) {
                GetunReadedmessage = c.unReadedmessage || 0
            }
        })
        if (GetunReadedmessage === 0) return
        var data = {mainuid, otheruid, GetunReadedmessage}
        var res = await this.MarkMsgsAsReaded(data)
        if (!res) return
        var {isMarked} = res
        if (isMarked) {
            this.contacts.forEach(c => {
                if (String(otheruid) == String(c?._id)) {
                    c.unReadedmessage = 0
                }
            })
        }
    },
    async GetUnreadedMsgList() {
        var res = await this.GetUnreadedMessageNum(this.MainUserData?._id)
        if (!res) return
        var {messages} = res
        this.contacts.forEach(user => {
            messages.forEach(msg => {
                if (String(msg.otherUserid) == String(user?._id)) {
                    user.unReadedmessage = Number(msg.numOfUnreadedMessages)
                }
            })
        })
    },
    async GetUsList() {
        this.contacts = []
        var glist = await this.GetUserFollowersFollowing
        this.contacts = glist
        if (this.contacts) {
            this.GetUnreadedMsgList()
        }
        this.updateOnlineList()
    },
    async selectUser(user) {
        this.messageBetweenUsers = []
        this.messagelistnum = 0
        this.selectedUser = user
        localStorage.setItem('chat_selected_user', JSON.stringify(user))
        await this.$nextTick()
        var firstuid = this.MainUserData?._id
        var seconduid = user._id
        var from = 0
        var ndata = {from, firstuid, seconduid}
        var res = await this.GetMsgsBetweebTwoUsers(ndata)
        if (!res) return
        var {msgs} = res
        if (!msgs) return
        this.messageBetweenUsers.push(...msgs)
        await this.$nextTick()
        this.scrollDownFunction()
        this.CallMarkMsgAsReaded(user)
    },
    async Sendmessage() {
        var content = this.messageToSend.text
        var sender = this.MainUserData?._id
        var receiver = this.selectedUser?._id
        var sdata = {content, sender, receiver}

        if (!this.uniqueOnlineUsers.includes(receiver)) {
            var res = await this.SendMessage(sdata)
            if (res) {
                this.messageBetweenUsers.push({content, sender, receiver})
                setTimeout(() => {
                    this.scrollDownFunction()
                }, 100)
            }
        } else {
            this.SendPrivateMessage(sdata).then(() => {
                this.messageBetweenUsers.push(sdata)   
            })
            setTimeout(() => {
                    this.scrollDownFunction()
                }, 100)
        }

        this.messageToSend.text = ''
        
    }
  }
}
</script>

<style scoped>
.chat-container {
    display: flex;
}

.chat-messages {
    flex: 1;
    padding: 10px;
}


</style>