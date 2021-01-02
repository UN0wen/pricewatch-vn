import { CookieWrapper } from "../utils/storage";
import { UserAuth } from "./context";

const currentUser = CookieWrapper.getCookie("userAuth");
const currentJWT = CookieWrapper.getCookie("jwt");
const user = currentUser
  ? currentUser
  : "";
const token = currentJWT
  ? currentJWT
  : "";
  
export const initialState : UserAuth  = {
  user: user || "",
  token: token || "",
  loading: false,
  errorMessage: null
};

export enum AuthReducerActions {
    REQUEST_LOGIN,
    LOGIN_SUCCESS,
    LOGOUT,
    LOGIN_ERROR
}

export interface AuthDispatch {
    type: AuthReducerActions,
    payload: { user: any; auth_token: string; },
    error: Error
}

export const AuthReducer = (initialState: any, action: AuthDispatch) => {
  switch (action.type) {
    case AuthReducerActions.REQUEST_LOGIN:
      return {
        ...initialState,
        loading: true
      };
    case AuthReducerActions.LOGIN_SUCCESS:
      return {
        ...initialState,
        user: action.payload.user,
        token: action.payload.auth_token,
        loading: false
      };
    case AuthReducerActions.LOGOUT:
      return {
        ...initialState,
        user: "",
        token: ""
      };
 
    case AuthReducerActions.LOGIN_ERROR:
      return {
        ...initialState,
        loading: false,
        errorMessage: action.error
      };
 
    default:
      throw new Error(`Unhandled action type: ${action.type}`);
  }
};