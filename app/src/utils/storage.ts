import Cookies from "universal-cookie";
const cookies = new Cookies();

export const CookieWrapper = (function () {
  const getCookie = function (): string {
    return cookies.get("userAuth");
  };

  const setCookie = function (userAuth: any, expire: Date) {
    cookies.set("userAuth", JSON.stringify(userAuth), {
      path: "/",
      expires: expire,
    });
  };

  const removeCookie = function () {
    cookies.remove("userAuth");
  };

  return {
    getCookie: getCookie,
    setCookie: setCookie,
    removeCookie: removeCookie,
  };
})();
