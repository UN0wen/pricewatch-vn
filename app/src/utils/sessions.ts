import { createContext } from 'react';
import Cookies from 'universal-cookie'

const cookies = new Cookies();

export var SessionContext = (function() {    
    var getJWT = function() : string {
      return cookies.get('jwt');   
    };
  
    var setJWT = function(jwt: string, expire: Date)  {
        cookies.set('jwt', jwt, { path: '/', expires: expire});
    };
    
    return {
      getJWT: getJWT,
      setJWT: setJWT,
    }
  
  })();
  
  export class UserCtx {
    username: string = ""
    jwt: string = ""
    auth: boolean = false
    constructor(username: string, jwt: string, auth: boolean) {
      this.username = username
      this.jwt = jwt
      this.auth = auth
    }
  }
  export const UserContext = createContext<UserCtx>(new UserCtx("", SessionContext.getJWT(), false))
