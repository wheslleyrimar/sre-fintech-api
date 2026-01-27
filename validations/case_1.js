import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    search: {
      executor: 'constant-vus',
      duration: '30s',
      vus: 25,
    },
  },
};

const base = __ENV.SRE_URL || 'http://localhost:8081';

export default function () {
  const term = 'checking';
  const url = `${base}/v1/search?term=${term}`;

  const res = http.get(url);

  check(res, { 'search status 200': (r) => r && r.status === 200 });
}
