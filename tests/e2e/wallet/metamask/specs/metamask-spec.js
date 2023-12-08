describe("Metamask Extension tests", () => {
  it("connect to DApp with Metamask extension", () => {
    cy.visit("/");
    cy.get("#connectButton").click();
    cy.acceptMetamaskAccess().should("be.true");
    cy.get("#connectButton").should("have.text", "Connected");
  });

  it("create transaction", () => {
    cy.visit("/");
    cy.get("#sendButton").click();
    cy.confirmMetamaskTransaction().should(($resp) => {
      expect($resp.confirmed).to.eq(true);
    });
  });
});
