import Cookies from "universal-cookie";
const cookies = new Cookies();

export const CookieWrapper = (function () {
  const getCookie = function (name:string) {
    return cookies.get(name);
  };

  const setCookie = function (name: string, data: any, expire: Date) {
    cookies.set(name, data, {
      path: "/",
      expires: expire,
    });
  };

  const removeCookie = function () {
    cookies.remove("userAuth");
    cookies.remove("jwt");
  };

  return {
    getCookie: getCookie,
    setCookie: setCookie,
    removeCookie: removeCookie,
  };
})();
