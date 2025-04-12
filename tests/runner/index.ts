import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';
import yaml from 'js-yaml';
import axios from 'axios';
import chalk from 'chalk';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

type TestCase = {
  name: string;
  baseUrl: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  url: string;
  headers?: Record<string, string>;
  body?: any;
  expect: {
    status: number;
    bodyContains?: string;
    bodyEquals?: any;
  };
};

const TEST_DIR = join(__dirname, 'yaml');

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

async function runTest(test: TestCase): Promise<{ name: string; passed: boolean; error?: string }> {
  try {
    const res = await axios({
      method: test.method,
      url: `${test.baseUrl}${test.url}`,
      headers: {
        Accept: 'application/json',
        ...(test.headers || {}),
      },
      data: test.body,
      validateStatus: () => true, // allow all status codes
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