import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';
import yaml from 'js-yaml';
import request from 'supertest';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

type TestCase = {
  name: string;
  baseUrl: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  url: string;
  body?: any;
  expect: {
    status: number;
    bodyContains?: string;
  };
};

const TEST_DIR = join(__dirname, 'yaml');

if (!fs.existsSync(TEST_DIR)) {
  console.error(`❌ Directory not found: ${TEST_DIR}`);
  process.exit(1);
}

async function runTest(test: TestCase) {
  try {
    const res = await request(test.baseUrl)
      [test.method.toLowerCase() as 'get' | 'post' | 'put' | 'delete'](test.url)
      .send(test.body || {})
      .set('Accept', 'application/json');

    if (res.status !== test.expect.status) {
      throw new Error(`${test.name} failed: expected status ${test.expect.status}, got ${res.status}`);
    }

    if (test.expect.bodyContains && !JSON.stringify(res.body).includes(test.expect.bodyContains)) {
      throw new Error(`${test.name} failed: response does not include expected content`);
    }

    console.log(`✅ ${test.name}`);
  } catch (err: any) {
    console.error(`❌ ${test.name} failed with error: ${err?.message || err}`);
    throw err;
  }
}

async function runAllTests() {
  const files = fs.readdirSync(TEST_DIR);

  for (const file of files) {
    if (!file.endsWith('.yaml')) continue;

    const filePath = join(TEST_DIR, file);
    console.log(`📄 Running tests from file: ${filePath}`);

    try {
      const doc = yaml.load(fs.readFileSync(filePath, 'utf8')) as TestCase[];

      if (!Array.isArray(doc)) {
        console.warn(`⚠️ Skipping invalid test file (not an array): ${file}`);
        continue;
      }

      for (const testCase of doc) {
        await runTest(testCase);
      }
    } catch (err: any) {
      console.error(`❌ Error while processing file "${file}":`, err?.message || JSON.stringify(err, null, 2));
      process.exit(1);
    }
  }

  console.log('🎉 All tests completed successfully');
}

runAllTests().catch((err: any) => {
  console.error('❌ Uncaught error during test execution:', err?.message || JSON.stringify(err, null, 2));
  process.exit(1);
});