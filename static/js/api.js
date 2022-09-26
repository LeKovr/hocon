var AppAPI = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

  // service.pb.ts
  var service_pb_exports = {};
  __export(service_pb_exports, {
    HoconService: () => HoconService,
    LampScene: () => LampScene
  });

  // ../fetch.pb.ts
  var b64 = new Array(64);
  var s64 = new Array(123);
  for (let i = 0; i < 64; )
    s64[b64[i] = i < 26 ? i + 65 : i < 52 ? i + 71 : i < 62 ? i - 4 : i - 59 | 43] = i++;
  function fetchReq(path, init) {
    const { pathPrefix, ...req } = init || {};
    const url = pathPrefix ? `${pathPrefix}${path}` : path;
    return fetch(url, req).then((r) => r.json().then((body) => {
      if (!r.ok) {
        throw body;
      }
      return body;
    }));
  }
  function isPlainObject(value) {
    const isObject = Object.prototype.toString.call(value).slice(8, -1) === "Object";
    const isObjLike = value !== null && isObject;
    if (!isObjLike || !isObject) {
      return false;
    }
    const proto = Object.getPrototypeOf(value);
    const hasObjectConstructor = typeof proto === "object" && proto.constructor === Object.prototype.constructor;
    return hasObjectConstructor;
  }
  function isPrimitive(value) {
    return ["string", "number", "boolean"].some((t) => typeof value === t);
  }
  function isZeroValuePrimitive(value) {
    return value === false || value === 0 || value === "";
  }
  function flattenRequestPayload(requestPayload, path = "") {
    return Object.keys(requestPayload).reduce(
      (acc, key) => {
        const value = requestPayload[key];
        const newPath = path ? [path, key].join(".") : key;
        const isNonEmptyPrimitiveArray = Array.isArray(value) && value.every((v) => isPrimitive(v)) && value.length > 0;
        const isNonZeroValuePrimitive = isPrimitive(value) && !isZeroValuePrimitive(value);
        let objectToMerge = {};
        if (isPlainObject(value)) {
          objectToMerge = flattenRequestPayload(value, newPath);
        } else if (isNonZeroValuePrimitive || isNonEmptyPrimitiveArray) {
          objectToMerge = { [newPath]: value };
        }
        return { ...acc, ...objectToMerge };
      },
      {}
    );
  }
  function renderURLSearchParams(requestPayload, urlPathParams = []) {
    const flattenedRequestPayload = flattenRequestPayload(requestPayload);
    const urlSearchParams = Object.keys(flattenedRequestPayload).reduce(
      (acc, key) => {
        const value = flattenedRequestPayload[key];
        if (urlPathParams.find((f) => f === key)) {
          return acc;
        }
        return Array.isArray(value) ? [...acc, ...value.map((m) => [key, m.toString()])] : acc = [...acc, [key, value.toString()]];
      },
      []
    );
    return new URLSearchParams(urlSearchParams).toString();
  }

  // service.pb.ts
  var LampScene = /* @__PURE__ */ ((LampScene2) => {
    LampScene2["UNKNOWN"] = "UNKNOWN";
    LampScene2["OFF"] = "OFF";
    LampScene2["NIGHT"] = "NIGHT";
    LampScene2["DAY"] = "DAY";
    return LampScene2;
  })(LampScene || {});
  var HoconService = class {
    static LampControl(req, initReq) {
      return fetchReq(`/api/lamp?${renderURLSearchParams(req, [])}`, { ...initReq, method: "GET" });
    }
  };
  return __toCommonJS(service_pb_exports);
})();
