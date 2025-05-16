import prompts from 'prompts';

async function confirm(question: string) {
  const answer = await prompts({
    type: 'confirm',
    name: 'confirmed',
    message: question,
    initial: false,
  });

  if (!answer.confirmed) {
    console.info('\r\n');
    console.info(`----ABORT SCRIPT----`);
    throw new Error(`Do not accept question: ${question}`);
  }
}

