import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.FileInputStream;
import java.util.List;
import java.util.Properties;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.route53.Route53Client;
import software.amazon.awssdk.services.route53.model.Change;
import software.amazon.awssdk.services.route53.model.ChangeAction;
import software.amazon.awssdk.services.route53.model.ChangeBatch;
import software.amazon.awssdk.services.route53.model.ChangeInfo;
import software.amazon.awssdk.services.route53.model.ChangeResourceRecordSetsRequest;
import software.amazon.awssdk.services.route53.model.ListResourceRecordSetsRequest;
import software.amazon.awssdk.services.route53.model.ResourceRecordSet;

public class TestChangeRecord {

  public static void main(String[] args) {

    Properties properties = new Properties();
    try {
      properties.load(new FileInputStream(ClassLoader.getSystemClassLoader().getResource("config").getPath()));
      System.out.println(properties.getProperty("ZoneName"));
    } catch (Exception e) {
    }

    Route53Client route53 = Route53Client.builder()
        .region(Region.AWS_GLOBAL)
        .build();

    List<ResourceRecordSet> resourceRecordSets = route53.listResourceRecordSets(
        ListResourceRecordSetsRequest.builder().hostedZoneId(properties.getProperty("ZoneId")).build()).resourceRecordSets();
    assertTrue(resourceRecordSets.size() > 0);
    ResourceRecordSet existingResourceRecordSet = resourceRecordSets.get(0);

    // Change Resource Record Sets
    ResourceRecordSet newResourceRecordSet = ResourceRecordSet.builder()
        .name(properties.getProperty("ZoneName"))
        .resourceRecords(existingResourceRecordSet.resourceRecords())
        .ttl(existingResourceRecordSet.ttl() + 100)
        .type(existingResourceRecordSet.type())
        .build();

    System.out.println(newResourceRecordSet.resourceRecords());
/*
    ChangeInfo changeInfo = route53.changeResourceRecordSets(ChangeResourceRecordSetsRequest.builder()
        .hostedZoneId(properties.getProperty("ZoneId"))
        .changeBatch(ChangeBatch.builder().comment("comment")
            .changes(/*Change.builder().action(
                        ChangeAction.DELETE)
                    .resourceRecordSet(
                        existingResourceRecordSet).build(),   //
                Change.builder().action(
                        ChangeAction.CREATE)
                    .resourceRecordSet(
                        newResourceRecordSet).build()).build()
        ).build()).changeInfo();
    assertValidChangeInfo(changeInfo);*/
  }

  /**
   * Asserts that the specified ChangeInfo is valid.
   *
   * @param change The ChangeInfo object to test.
   */
  private static void assertValidChangeInfo(ChangeInfo change) {
    assertNotNull(change.id());
    assertNotNull(change.status());
    assertNotNull(change.submittedAt());
  }
}
