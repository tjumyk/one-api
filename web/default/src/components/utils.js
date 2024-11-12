import { API, showError } from '../helpers';

export async function getOAuthState() {
  const res = await API.get('/api/oauth/state');
  const { success, message, data } = res.data;
  if (success) {
    return data;
  } else {
    showError(message);
    return '';
  }
}

export async function onGitHubOAuthClicked(github_client_id) {
  const state = await getOAuthState();
  if (!state) return;
  let redirect_uri = `${window.location.origin}/oauth-callback`;
  let url = `https://id.lmzgc.cn:8000/oauth/connect?client_id=${github_client_id}&redirect_url=${redirect_uri}&state=${state}`;
  window.open(url);
}

export async function onLarkOAuthClicked(lark_client_id) {
  const state = await getOAuthState();
  if (!state) return;
  let redirect_uri = `${window.location.origin}/oauth/lark`;
  window.open(
    `https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=${redirect_uri}&app_id=${lark_client_id}&state=${state}`
  );
}