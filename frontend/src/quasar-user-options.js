
import './styles/quasar.sass'
import iconSet from 'quasar/icon-set/eva-icons.js'
import '@quasar/extras/eva-icons/eva-icons.css'
import { Notify } from 'quasar'

// To be used on app.use(Quasar, { ... })
export default {
  config: {
    notify: {}
  },
  plugins: {
    Notify
  },
  iconSet: iconSet
}