class Key extends MUDObject {
    private key: string = ""

    Matches(id: string) {
        return this.key == id;
    }
}