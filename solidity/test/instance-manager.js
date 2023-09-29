const InstanceManager = artifacts.require("InstanceManager");

const instantiationProof = JSON.parse(`{"a":["4536485193013306460708422378779740504073699311711912164105515306662962106160","20478655771368444655258849833648885959074986187725633342309936581101655763199"],"b":[["16188466433589100578671328381999457373493483661329235890095974074238720872147","13590391848409948134818716802703452199484405011458992824460642626254099279904"],["13908679637417014569172588734800053987335647950953482063203361693744145343606","16496432789690918117327257320320601287104757862848192579653697423588139775493"]],"c":["16219772369300257498826416241509833475061364010583533865874537611515148161072","20577305547706293897880892356203810791683567296948031334849213594489649744245"],"publicInput":["19092878713960129258038638675631701634919395924240493996892111696307011520512","17529240440588388804137033526724658072525540712954021729923811435556868282112"]}`)

contract('InstanceManager', (accounts) => {
  it('should put 10000 MetaCoin in the first account', async () => {
    const metaCoinInstance = await MetaCoin.deployed();
    const balance = await metaCoinInstance.getBalance.call(accounts[0]);

    assert.equal(balance.valueOf(), 10000, "10000 wasn't in the first account");
  });
  it('should call a function that depends on a linked library', async () => {
    const metaCoinInstance = await MetaCoin.deployed();
    const metaCoinBalance = (await metaCoinInstance.getBalance.call(accounts[0])).toNumber();
    const metaCoinEthBalance = (await metaCoinInstance.getBalanceInEth.call(accounts[0])).toNumber();

    assert.equal(metaCoinEthBalance, 2 * metaCoinBalance, 'Library function returned unexpected function, linkage may be broken');
  });
  it('should send coin correctly', async () => {
    const metaCoinInstance = await MetaCoin.deployed();

    // Setup 2 accounts.
    const accountOne = accounts[0];
    const accountTwo = accounts[1];

    // Get initial balances of first and second account.
    const accountOneStartingBalance = (await metaCoinInstance.getBalance.call(accountOne)).toNumber();
    const accountTwoStartingBalance = (await metaCoinInstance.getBalance.call(accountTwo)).toNumber();

    // Make transaction from first account to second.
    const amount = 10;
    await metaCoinInstance.sendCoin(accountTwo, amount, { from: accountOne });

    // Get balances of first and second account after the transactions.
    const accountOneEndingBalance = (await metaCoinInstance.getBalance.call(accountOne)).toNumber();
    const accountTwoEndingBalance = (await metaCoinInstance.getBalance.call(accountTwo)).toNumber();

    assert.equal(accountOneEndingBalance, accountOneStartingBalance - amount, "Amount wasn't correctly taken from the sender");
    assert.equal(accountTwoEndingBalance, accountTwoStartingBalance + amount, "Amount wasn't correctly sent to the receiver");
  });
});
