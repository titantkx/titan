import fs from 'fs';
import prompts from 'prompts';

function writeDeployInfo(infos = {}) {
  fs.rmSync('deploy_infos', { recursive: true, force: true });
  fs.mkdirSync('deploy_infos');
  fs.writeFileSync(
    `deploy_infos/deploy_infos.json`,
    JSON.stringify(
      infos,
      null,
      2 //eslint-disable-line no-magic-numbers
    )
  );
}

function writeCSVData(filename: string, infos: Array<string>[] = []) {
  fs.rmSync(`csv_data/${filename}.csv`, { recursive: true, force: true });
  if (!fs.existsSync('csv_data')) {
    fs.mkdirSync('csv_data');
  }
  var text2Write = '';
  infos.forEach((c) => {
    text2Write += c.join(',') + '\r\n';
  });
  fs.writeFileSync(`csv_data/${filename}.csv`, text2Write);
}

function writeGasStatisticInfos(filename: string, infos = {}) {
  fs.rmSync(`gas_statistic_infos/${filename}.json`, { recursive: true, force: true });
  if (!fs.existsSync('gas_statistic_infos')) {
    fs.mkdirSync('gas_statistic_infos');
  }
  fs.writeFileSync(
    `gas_statistic_infos/${filename}.json`,
    JSON.stringify(
      infos,
      null,
      2 //eslint-disable-line no-magic-numbers
    )
  );
}

function writeGasStatisticCSVInfos(filename: string, infos: Array<string>[] = []) {
  fs.rmSync(`gas_statistic_infos/${filename}.csv`, { recursive: true, force: true });
  if (!fs.existsSync('gas_statistic_infos')) {
    fs.mkdirSync('gas_statistic_infos');
  }
  var text2Write = '';
  infos.forEach((c) => {
    text2Write += c.join(',') + '\r\n';
  });
  fs.writeFileSync(`gas_statistic_infos/${filename}.csv`, text2Write);
}

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

export default {
  writeDeployInfo,
  writeCSVData,
  writeGasStatisticInfos,
  writeGasStatisticCSVInfos,
  confirm,
};
