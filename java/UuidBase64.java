import java.nio.ByteBuffer;
import java.util.UUID;
import java.util.Base64;
import java.io.UnsupportedEncodingException;

class UuidBase64 {
    public static byte[] getBytesFromUUID(UUID uuid) {
        ByteBuffer bb = ByteBuffer.wrap(new byte[16]);
        bb.putLong(uuid.getMostSignificantBits());
        bb.putLong(uuid.getLeastSignificantBits());

        return bb.array();
    }

    public static UUID getUUIDFromBytes(byte[] bytes) {
        ByteBuffer byteBuffer = ByteBuffer.wrap(bytes);
        return new UUID(byteBuffer.getLong(), byteBuffer.getLong());
    }

    public static UUID getUUIDFromBase64(String data) throws UnsupportedEncodingException {
        byte[] decodedString = Base64.getDecoder().decode(data.getBytes("UTF-8"));
        return UuidBase64.getUUIDFromBytes(decodedString);
    }

    public static String makeBase64Uuid(UUID u) {
        return new String(Base64.getEncoder().encode(UuidBase64.getBytesFromUUID(u)));
    }
}
