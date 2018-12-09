import junit.framework.Assert;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.transaction.TransactionConfiguration;
import tds.student.performance.dao.utils.UuidAdapter;

import java.util.UUID;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration("/your-test-context.xml")
@TransactionConfiguration
public class UuidAdapterTest {

    @Test
    public void should_Get_16_Bytes_From_a_UUID() {
        UUID uuid = UUID.randomUUID();

        byte[] result = UuidAdapter.getBytesFromUUID(uuid);

        Assert.assertEquals("Expected result to be a byte array w/16 elements.", 16, result.length);
    }

    @Test
    public void should_Reconstruct_Same_UUID_From_Byte_Array() {
        UUID uuid = UUID.randomUUID();

        byte[] bytes = UuidAdapter.getBytesFromUUID(uuid);
        UUID reconstructedUuid = UuidAdapter.getUUIDFromBytes(bytes);

        Assert.assertEquals(uuid, reconstructedUuid);
    }

    @Test
    public void should_Not_Generate_the_Same_UUID_From_Bytes() {
        UUID uuid = UUID.fromString("9f881758-0b4a-4eaa-b59f-b6dea0934223");

        byte[] result = UuidAdapter.getBytesFromUUID(uuid);
        UUID newUuid = UUID.nameUUIDFromBytes(result);

        Assert.assertFalse(uuid.equals(newUuid));
    }
}