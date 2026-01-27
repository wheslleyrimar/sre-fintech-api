import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    report: {
      executor: 'constant-vus',
      duration: '30s',
      vus: 10,
    },
  },
};

const base = __ENV.SRE_URL || 'http://localhost:8081';

export default function () {
  const url = `${base}/v1/report`;

  const res = http.get(url);

  check(res, { 'report status 200': (r) => r && r.status === 200 });
}
