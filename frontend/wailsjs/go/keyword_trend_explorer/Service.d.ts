// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {resp} from '../models';
import {pagination} from '../models';
import {keyword_trend_explorer} from '../models';

export function DeleteKeywordTrendExplorer(arg1:number):Promise<resp.Response>;

export function GetKeywordTrendExplorer(arg1:number):Promise<resp.Response>;

export function GetKeywordTrendExplorers(arg1:pagination.Pagination):Promise<resp.Response>;

export function GetKeywordsByExplorers(arg1:number,arg2:pagination.Pagination):Promise<resp.Response>;

export function KeywordTrendExplorerRequest(arg1:keyword_trend_explorer.KeywordTrendExplorerRequestPayload):Promise<resp.Response>;