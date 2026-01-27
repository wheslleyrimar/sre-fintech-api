import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  scenarios: {
    tariff_adjustment: {
      executor: 'constant-vus',
      duration: '60s',
      vus: 1,
    },
  },
};

const base = __ENV.SRE_URL || 'http://localhost:8081';

const accountIds = [
  '23423da1-1570-4c3f-8384-2e0aa022e486',
  '1f163dcb-92fe-4503-b5bb-95c6d518e7af',
  '4c835201-99eb-4330-84f2-a8f689636f38',
  '68606754-4340-482d-94c7-fabe890a901c',
  'a75602d8-a451-4d7b-9075-c8e4929d9332',
  '02b157f5-8b2b-4391-adfe-7f918635ca05',
  '99fadf7c-279e-476b-b651-061ac0709c80',
  '3ebaed0a-7b67-4987-8fa3-ed7ebb7290e2',
  '80e8ae8b-7d16-4851-a971-27c477e6d754',
  'ec358bca-a616-46ba-8029-11281af697d0',
  '0dd2b43e-7586-4a7c-807e-096d96162d73',
  '1510f6e3-c340-4eac-be1f-161fe95d2567',
  '67ad45ce-e5e1-4c5d-83ce-56627fcf9584',
  'e967498a-128e-4db8-a560-4991efee833e',
  '6112f772-a514-459e-90e9-8ba1434d7017',
  '9a309334-19ef-43cc-bbda-d0700bfe2a9b',
  'e0e806a5-d58b-491f-8026-e20099f23f0e',
  '8f75d6ee-c8d7-4819-a23b-0aa29059f0a5',
  'e348dc40-598a-4dc2-ac93-00afe68d053b',
  '739efbd3-1d77-4364-9680-5f8563f60ae6',
  '23e59477-4b6b-4c73-9b1a-38599478ff09',
  '746ba36d-7f81-4a92-a75d-dcd7ce5af393',
  '4eea9b6c-4160-45d6-88fe-fb8678b7096e',
  '96f54994-b4c4-4084-a058-1ce1e7061436',
  'c07ca196-1928-4c65-a7f2-395f90ecfd1b',
  'b60ec4af-5c14-4b51-8fb8-a29f8af0bd81',
  '82daf290-f9f0-475a-9dd2-30da9d997960',
  '2c2e52f3-4bea-45d9-9f3a-a41093093f0a',
  'ac42d9fd-286d-437e-b0de-b4e7434dc832',
  'e77dd800-3a4d-4277-8b24-d275980bc931',
  '2ee7b40b-4dae-4867-a010-b1d7b77068ea',
  '32204424-da1c-4704-adb9-595e636550a9',
  'a71a0b8f-4907-42d2-899d-d1195511f301',
  '32163177-663c-4241-83e7-553b76793f94',
  '08cef8e1-a232-45b7-95b1-366943645e81',
  'ccefab25-447a-4331-b077-bc6599ebeca4',
  '8ccb4d81-6499-4cd0-aaf0-2b431027b78a',
];

export function setup() {
  const index = (parseInt(__ENV.GROUP_NUMBER, 10) || 1) - 1;
  const id = accountIds[index % accountIds.length];
  return { id };
}

export default function (data) {
  const adjustUrl = `${base}/v1/accounts/${data.id}/tariff-adjustments`;
  const accountUrl = `${base}/v1/accounts/${data.id}`;

  const newFee = randomIntBetween(10, 200) / 10;
  const payload = JSON.stringify({ new_fee: newFee });

  const params = { headers: { 'Content-Type': 'application/json' } };

  const postRes = http.post(adjustUrl, payload, params);
  if (postRes.status !== 204 && postRes.status !== 200) {
    return;
  }

  sleep(2);

  const getRes = http.get(accountUrl);
  check(getRes, {
    'account monthly_fee matches requested new_fee': (r) =>
      r.status === 200 && Number(r.json().monthly_fee) === newFee,
  });
}
