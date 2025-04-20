import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';
import yaml from 'js-yaml';
import axios from 'axios';
import chalk from 'chalk';
import request from 'supertest';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

type TestCase = {
  name: string;
  baseUrl: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  url: string;
  auth?: boolean;
  headers?: Record<string, string>;
  body?: any;
  expect: {
    status: number;
    bodyContains?: string;
    bodyEquals?: any;
  };
};

const TEST_DIR = join(__dirname, 'yaml');
const GATEWAY_BASE_URL = 'http://api-gateway:3000';
const TEST_USER = { username: 'testuser', password: 'password123' };

let jwtToken: string = '';

function getAllYamlFiles(dir: string): string[] {
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  return entries.flatMap((entry) => {
    const fullPath = join(dir, entry.name);
    if (entry.isDirectory()) {
      return getAllYamlFiles(fullPath);
    }
    if (entry.isFile() && fullPath.endsWith('.yaml')) {
      return [fullPath];
    }
    return [];
  });
}

async function setupTestUser(): Promise<void> {
  try {
    // ユーザー登録（存在すればエラーだが無視）
    await axios.post(`${GATEWAY_BASE_URL}/auth/register`, TEST_USER);
    console.log('✅ Test user registered');
  } catch (err: any) {
    if (err.response?.status === 409) {
      console.log('⚠️ Test user already exists (409)');
    } else {
      console.error('❌ Failed to register test user:', err.message);
    }
  }

  // トークン取得
  try {
    const res = await axios.post(`${GATEWAY_BASE_URL}/auth/login`, TEST_USER);
    jwtToken = res.data.access_token;
    console.log('✅ JWT token obtained:');
  } catch (err: any) {
    console.error('❌ Failed to get JWT token:', err.message);
    process.exit(1);
  }
}

async function runTest(test: TestCase): Promise<{ name: string; passed: boolean; error?: string }> {
  try {
    const headers: Record<string, string> = {
      Accept: 'application/json',
      ...(test.headers || {}),
    };

    const isGatewayRequest = test.baseUrl.includes('api-gateway');
    const isProtected = test.auth !== false;
    const hasToken = !!jwtToken;

    console.log(chalk.yellow(`🧐 Header Injection Check for "${test.name}"`));
    console.log(chalk.gray(`  isGatewayRequest: ${isGatewayRequest}`));
    console.log(chalk.gray(`  isProtected     : ${isProtected}`));
    console.log(chalk.gray(`  hasToken        : ${hasToken}`));

    if (isGatewayRequest && isProtected && jwtToken) {
      headers['Authorization'] = `Bearer ${jwtToken}`;
      console.log(chalk.gray(`🔐 Added Authorization header for "${test.name}"`));
    }

    console.log(headers);

    const res = await axios({
      method: test.method,
      url: `${test.baseUrl}${test.url}`,
      headers,
      data: test.body,
      validateStatus: () => true,
    });

    if (res.status !== test.expect.status) {
      return {
        name: test.name,
        passed: false,
        error: `expected status ${test.expect.status}, got ${res.status}`,
      };
    }

    if (test.expect.bodyContains && !JSON.stringify(res.data).includes(test.expect.bodyContains)) {
      return {
        name: test.name,
        passed: false,
        error: `expected body to include "${test.expect.bodyContains}"`,
      };
    }

    if (test.expect.bodyEquals && JSON.stringify(res.data) !== JSON.stringify(test.expect.bodyEquals)) {
      return {
        name: test.name,
        passed: false,
        error: `expected exact match in body`,
      };
    }

    return { name: test.name, passed: true };
  } catch (err: any) {
    return {
      name: test.name,
      passed: false,
      error: err.message || 'Unknown error',
    };
  }
}

async function runAllTests() {
  // トークンは最初に取得して使いまわす
  await setupTestUser();

  const files = getAllYamlFiles(TEST_DIR);
  console.log(chalk.blueBright('📂 Test files found:'), files);

  const results: { name: string; passed: boolean; error?: string }[] = [];

  for (const file of files) {
    try {
      const raw = fs.readFileSync(file, 'utf8');
      const doc = yaml.load(raw) as TestCase[];

      if (!Array.isArray(doc)) {
        console.warn(chalk.yellow(`⚠️  Skipping invalid test file: ${file}`));
        continue;
      }

      for (const testCase of doc) {
        await cleanupServices();

        const result = await runTest(testCase);
        results.push(result);

        if (result.passed) {
          console.log(chalk.green(`✅ ${testCase.name}`));
        } else {
          console.log(chalk.red(`❌ ${testCase.name} - ${result.error}`));
        }
      }
    } catch (err: any) {
      console.error(chalk.red(`❌ Error loading file ${file}: ${err.message || err}`));
    }
  }

  const failed = results.filter((r) => !r.passed);

  console.log();
  console.log(chalk.blue('📋 Test Results Summary:'));
  console.log(chalk.green(`  ✅ Passed: ${results.length - failed.length}`));
  console.log(chalk.red(`  ❌ Failed: ${failed.length}`));

  if (failed.length === 0) {
    console.log(chalk.greenBright('\n🎉 All tests passed!'));
  } else {
    console.log(chalk.redBright('\n❌ Some tests failed!'));
  }
}

runAllTests().catch((err) => {
  console.error('❌ Test runner crashed:', err.message);
  process.exit(1);
});

async function cleanupServices() {
  const services = [
    { name: 'auth-service', url: 'http://auth-service:8000' },
    { name: 'post-service', url: 'http://post-service:3000' },
    { name: 'emotion-service', url: 'http://emotion-service:8080' },
    { name: 'tag-service', url: 'http://tag-service:4000' }
  ];

  for (const service of services) {
    try {
      await request(service.url)
        .post('/test/cleanup')
        .expect(200);
      console.log(`✅ ${service.name} cleanup done`);
    } catch (e) {
      console.warn(`⚠️ ${service.name} cleanup failed:`, (e as Error).message);
    }
  }
}