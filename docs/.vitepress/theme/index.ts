import DefaultTheme from 'vitepress/theme'
import CustomLayout from './Layout.vue'
import './style.css'

export default {
  ...DefaultTheme,
  Layout: CustomLayout
}
