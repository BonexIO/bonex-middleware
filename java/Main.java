import java.nio.ByteBuffer;
import java.util.UUID;
import java.util.Base64;
import java.io.UnsupportedEncodingException;

class Main {
    public static void main(String[] args) {
        UUID uuid = UUID.randomUUID();
        System.out.println("Random   UUID: " + uuid.toString());
        String result = UuidBase64.makeBase64Uuid(uuid);
        System.out.println("Generated Base64: " + result);

        try {
            String data = new String(result);
            UUID u = UuidBase64.getUUIDFromBase64(data);
            System.out.println("Restored UUID: " + u.toString());
        } catch ( UnsupportedEncodingException e) {
            System.out.println("Unsupported character set");
        }
    }
}