import store from "@/store";

export default function NumAuthGuard(to, from, next) {
  if (store.state.Auth.authData) {
    if (store.getters.isLoggedIn) {
      next('/');
    }
  } else {
    next();
  }
}