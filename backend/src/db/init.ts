import { createClient } from '@libsql/client';
import { readFileSync } from 'fs';
import { join } from 'path';
import { dirname } from 'path';

const __dirname = dirname(process.argv[1]);

async function initDatabase() {
  try {
    const db = createClient({
      url: process.env.TURSO_CONNECTION_URL || 'file:./dev.db',
      authToken: process.env.TURSO_AUTH_TOKEN
    });

    console.log('Initializing database...');

    // Read schema file
    const schemaPath = join(__dirname, 'schema.sql');
    const schema = readFileSync(schemaPath, 'utf-8');

    // Split by semicolons and execute each statement
    const statements = schema
      .split(';')
      .map(stmt => stmt.trim())
      .filter(stmt => stmt.length > 0);

    for (const statement of statements) {
      console.log(`Executing: ${statement.substring(0, 50)}...`);
      await db.execute(statement);
    }

    console.log('✓ Database initialized successfully');
    process.exit(0);
  } catch (error) {
    console.error('✗ Failed to initialize database:', error);
    process.exit(1);
  }
}

initDatabase();
