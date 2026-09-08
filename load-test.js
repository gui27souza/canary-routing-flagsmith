import http from 'k6/http'
import { check, sleep } from 'k6'

    export const options = {

      stages: [
        { duration: '5s', target: 20 },
        { duration: '15s', target: 50 },
        { duration: '5s', target: 0 },
      ],

      thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<15', 'p(99)<30'],
      },
    }
    
    export default function () {
      const url = 'http://127.0.0.1:8080/decide'

      const randomId = Math.floor(Math.random() * 1000000)
      const payload = JSON.stringify({
        user_id: `user_${randomId}`,
        email: `user_${randomId}@platform.test`,
        country: 'BR',
      })

      const params = {
        headers: {
          'Content-Type': 'application/json',
        },
      }

      const res = http.post(url, payload, params)

      check(res, {
        'status is 200': (r) => r.status === 200,
        'has routing decision': (r) => r.body && r.body.includes('target'),
      })

      sleep(0.01)
    }