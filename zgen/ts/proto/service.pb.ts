/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"

export enum LampScene {
  UNKNOWN = "UNKNOWN",
  OFF = "OFF",
  NIGHT = "NIGHT",
  DAY = "DAY",
}

export type LampStatus = {
  id?: string
  scene?: LampScene
}

export class HoconService {
  static LampControl(req: LampStatus, initReq?: fm.InitReq): Promise<LampStatus> {
    return fm.fetchReq<LampStatus, LampStatus>(`/api/lamp?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"})
  }
}