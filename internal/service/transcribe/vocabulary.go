package transcribe

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func ResourceVocabulary() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVocabularyCreate,
		ReadWithoutTimeout:   resourceVocabularyRead,
		UpdateWithoutTimeout: resourceVocabularyUpdate,
		DeleteWithoutTimeout: resourceVocabularyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"download_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"language_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(validateLanguageCodes(types.LanguageCode("").Values()), false),
			},
			"phrases": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     256,
				ExactlyOneOf: []string{"phrases", "vocabulary_file_uri"},
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
			"vocabulary_file_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"phrases", "vocabulary_file_uri"},
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},
			"vocabulary_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

const (
	ResNameVocabulary = "transcribe"
)

func resourceVocabularyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).TranscribeConn

	in := &transcribe.CreateVocabularyInput{
		VocabularyName: aws.String(d.Get("vocabulary_name").(string)),
		LanguageCode:   types.LanguageCode(d.Get("language_code").(string)),
	}

	if v, ok := d.GetOk("vocabulary_file_uri"); ok {
		in.VocabularyFileUri = aws.String(v.(string))
	}

	if v, ok := d.GetOk("phrases"); ok {
		in.Phrases = expandPhrases(v.([]interface{}))
	}

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	if len(tags) > 0 {
		in.Tags = Tags(tags.IgnoreAWS())
	}

	out, err := conn.CreateVocabulary(ctx, in)
	if err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionCreating, ResNameVocabulary, d.Get("vocabulary_name").(string), err)
	}

	if out == nil {
		return names.DiagError(names.Transcribe, names.ErrActionCreating, ResNameVocabulary, d.Get("vocabulary_name").(string), errors.New("empty output"))
	}

	d.SetId(aws.ToString(out.VocabularyName))

	if _, err := waitVocabularyCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionWaitingForCreation, ResNameVocabulary, d.Id(), err)
	}

	return resourceVocabularyRead(ctx, d, meta)
}

func resourceVocabularyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).TranscribeConn

	out, err := FindVocabularyByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Transcribe Vocabulary (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionReading, ResNameVocabulary, d.Id(), err)
	}

	arn := arn.ARN{
		AccountID: meta.(*conns.AWSClient).AccountID,
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "transcribe",
		Region:    meta.(*conns.AWSClient).Region,
		Resource:  fmt.Sprintf("vocabulary/%s", d.Id()),
	}.String()

	d.Set("arn", arn)
	d.Set("download_uri", out.DownloadUri)
	d.Set("vocabulary_name", out.VocabularyName)
	d.Set("language_code", out.LanguageCode)

	tags, err := ListTags(ctx, conn, arn)
	if err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionReading, ResNameVocabulary, d.Id(), err)
	}

	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionSetting, ResNameVocabulary, d.Id(), err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionSetting, ResNameVocabulary, d.Id(), err)
	}

	return nil
}

func resourceVocabularyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).TranscribeConn

	if d.HasChangesExcept("tags", "tags_all") {
		in := &transcribe.UpdateVocabularyInput{
			VocabularyName: aws.String(d.Id()),
			LanguageCode:   types.LanguageCode(d.Get("language_code").(string)),
		}

		if d.HasChanges("vocabulary_file_uri", "phrases") {
			if d.Get("vocabulary_file_uri").(string) != "" {
				in.VocabularyFileUri = aws.String(d.Get("vocabulary_file_uri").(string))
			} else {
				in.Phrases = expandPhrases(d.Get("phrases").([]interface{}))
			}
		}

		log.Printf("[DEBUG] Updating Transcribe Vocabulary (%s): %#v", d.Id(), in)
		_, err := conn.UpdateVocabulary(ctx, in)
		if err != nil {
			return names.DiagError(names.Transcribe, names.ErrActionUpdating, ResNameVocabulary, d.Id(), err)
		}

		if _, err := waitVocabularyUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return names.DiagError(names.Transcribe, names.ErrActionWaitingForUpdate, ResNameVocabulary, d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(ctx, conn, d.Get("arn").(string), o, n); err != nil {
			return diag.Errorf("error updating Transcribe Vocabulary (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceVocabularyRead(ctx, d, meta)
}

func resourceVocabularyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).TranscribeConn

	log.Printf("[INFO] Deleting Transcribe Vocabulary %s", d.Id())

	_, err := conn.DeleteVocabulary(ctx, &transcribe.DeleteVocabularyInput{
		VocabularyName: aws.String(d.Id()),
	})

	var badRequestException *types.BadRequestException
	if errors.As(err, &badRequestException) {
		return nil
	}

	if err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionDeleting, ResNameVocabulary, d.Id(), err)
	}

	if _, err := waitVocabularyDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return names.DiagError(names.Transcribe, names.ErrActionWaitingForDeletion, ResNameVocabulary, d.Id(), err)
	}

	return nil
}

func waitVocabularyCreated(ctx context.Context, conn *transcribe.Client, id string, timeout time.Duration) (*transcribe.GetVocabularyOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   vocabularyStatus(types.VocabularyStatePending),
		Target:                    vocabularyStatus(types.VocabularyStateReady),
		Refresh:                   statusVocabulary(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
		Delay:                     30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*transcribe.GetVocabularyOutput); ok {
		if status := out.VocabularyState; status == types.VocabularyStateFailed {
			return out, errors.New(aws.ToString(out.FailureReason))
		}
		return out, err
	}

	return nil, err
}

func waitVocabularyUpdated(ctx context.Context, conn *transcribe.Client, id string, timeout time.Duration) (*transcribe.GetVocabularyOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   vocabularyStatus(types.VocabularyStatePending),
		Target:                    vocabularyStatus(types.VocabularyStateReady),
		Refresh:                   statusVocabulary(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
		Delay:                     30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*transcribe.GetVocabularyOutput); ok {
		if status := out.VocabularyState; status == types.VocabularyStateFailed {
			return out, errors.New(aws.ToString(out.FailureReason))
		}
		return out, err
	}

	return nil, err
}

func waitVocabularyDeleted(ctx context.Context, conn *transcribe.Client, id string, timeout time.Duration) (*transcribe.GetVocabularyOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: vocabularyStatus(types.VocabularyStatePending),
		Target:  []string{},
		Refresh: statusVocabulary(ctx, conn, id),
		Timeout: timeout,
		Delay:   30 * time.Second,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*transcribe.GetVocabularyOutput); ok {
		if status := out.VocabularyState; status == types.VocabularyStateFailed {
			return out, errors.New(aws.ToString(out.FailureReason))
		}
		return out, err
	}

	return nil, err
}

func statusVocabulary(ctx context.Context, conn *transcribe.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := FindVocabularyByName(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.VocabularyState), nil
	}
}

func FindVocabularyByName(ctx context.Context, conn *transcribe.Client, id string) (*transcribe.GetVocabularyOutput, error) {
	in := &transcribe.GetVocabularyInput{
		VocabularyName: aws.String(id),
	}

	out, err := conn.GetVocabulary(ctx, in)

	var badRequestException *types.BadRequestException
	if errors.As(err, &badRequestException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: in,
		}
	}

	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func vocabularyStatus(in ...types.VocabularyState) []string {
	var s []string

	for _, v := range in {
		s = append(s, string(v))
	}

	return s
}

func expandPhrases(in []interface{}) []string {
	var out []string

	for _, val := range in {
		out = append(out, val.(string))
	}
	return out
}
