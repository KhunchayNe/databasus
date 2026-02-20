import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { RestoreTarget } from '../model/RestoreTarget';
import type { CreateRestoreTargetWithDBRequest, UpdateRestoreTargetRequest } from '../model/CreateRestoreTargetRequest';

export const restoreTargetApi = {
  async getRestoreTargets(): Promise<RestoreTarget[]> {
    const requestOptions: RequestOptions = new RequestOptions();
    const url = `${getApplicationServer()}/api/v1/restore-targets`;
    return apiHelper.fetchGetJson(url, requestOptions);
  },

  async getRestoreTarget(id: string): Promise<RestoreTarget> {
    const requestOptions: RequestOptions = new RequestOptions();
    const url = `${getApplicationServer()}/api/v1/restore-targets/${id}`;
    return apiHelper.fetchGetJson(url, requestOptions);
  },

  async createRestoreTarget(request: CreateRestoreTargetWithDBRequest): Promise<RestoreTarget> {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(request));
    const url = `${getApplicationServer()}/api/v1/restore-targets`;
    return apiHelper.fetchPostJson(url, requestOptions);
  },

  async updateRestoreTarget(
    id: string,
    request: UpdateRestoreTargetRequest,
  ): Promise<RestoreTarget> {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(request));
    const url = `${getApplicationServer()}/api/v1/restore-targets/${id}`;
    return apiHelper.fetchPutJson(url, requestOptions);
  },

  async deleteRestoreTarget(id: string): Promise<{ message: string }> {
    const requestOptions: RequestOptions = new RequestOptions();
    const url = `${getApplicationServer()}/api/v1/restore-targets/${id}`;
    return apiHelper.fetchDeleteJson(url, requestOptions);
  },
};
